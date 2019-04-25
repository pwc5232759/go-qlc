package dpos

import (
	"github.com/bluele/gcache"
	"github.com/qlcchain/go-qlc/consensus"
	"github.com/qlcchain/go-qlc/ledger/process"
	"sync"
	"time"

	"github.com/qlcchain/go-qlc/common"
	"github.com/qlcchain/go-qlc/common/event"
	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/config"
	"github.com/qlcchain/go-qlc/ledger"
	"github.com/qlcchain/go-qlc/log"
	"github.com/qlcchain/go-qlc/p2p"
	"github.com/qlcchain/go-qlc/p2p/protos"
	"go.uber.org/zap"
)

const (
	repTimeout                        = 5 * time.Minute
	msgCacheSize                      = 65536
	uncheckedTimeout                  = 5 * time.Minute
	voteCacheTimeout                  = 5 * time.Minute
	refreshPriInterval                = 5 * time.Minute
	findOnlineRepresentativesInterval = 2 * time.Minute
)

var (
	localRepAccount sync.Map
)

type DPoS struct {
	ledger         *ledger.Ledger
	acTrx          *ActiveTrx
	accounts       []*types.Account
	onlineReps     sync.Map
	logger         *zap.SugaredLogger
	cfg            *config.Config
	eb             event.EventBus
	lv             *process.LedgerVerifier
	uncheckedCache gcache.Cache //gap blocks
	voteCache      gcache.Cache //vote blocks
	quitCh         chan bool
}

func NewDPoS(cfg *config.Config, accounts []*types.Account, eb event.EventBus) *DPoS {
	acTrx := NewActiveTrx()
	l := ledger.NewLedger(cfg.LedgerDir(), eb)

	dps := &DPoS{
		ledger:         l,
		acTrx:          acTrx,
		accounts:       accounts,
		logger:         log.NewLogger("dpos"),
		cfg:            cfg,
		eb:             eb,
		lv:             process.NewLedgerVerifier(l),
		uncheckedCache: gcache.New(msgCacheSize).LRU().Expiration(uncheckedTimeout).Build(),
		voteCache:      gcache.New(msgCacheSize).LRU().Expiration(voteCacheTimeout).Build(),
		quitCh:			make(chan bool, 1),
	}

	dps.acTrx.SetDposService(dps)
	return dps
}

func (dps *DPoS) Init() {
	if len(dps.accounts) != 0 {
		dps.refreshAccount()
	}
}

func (dps *DPoS) Start() {
	dps.logger.Info("DPOS service started!")

	dps.acTrx.start()
	timerFindOnlineRep := time.NewTicker(findOnlineRepresentativesInterval)
	timerRefreshPri := time.NewTicker(refreshPriInterval)

	for {
		select {
		case <-dps.quitCh:
			dps.logger.Info("Stopped DPOS.")
			return
		case <-timerRefreshPri.C:
			dps.logger.Info("refresh pri info.")
			go dps.refreshAccount()
		case <-timerFindOnlineRep.C:
			dps.logger.Info("begin Find Online Representatives.")
			go func() {
				err := dps.findOnlineRepresentatives()
				if err != nil {
					dps.logger.Error(err)
				}
				dps.cleanOnlineReps()
			}()
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (dps *DPoS) Stop() {
	dps.logger.Info("DPOS service stopped!")
	dps.acTrx.stop()
	dps.quitCh <- true
}

func (dps *DPoS) enqueueUnchecked(hash types.Hash, depHash types.Hash, bs *consensus.BlockSource) {
	if !dps.uncheckedCache.Has(depHash) {
		blocks := new(sync.Map)
		blocks.Store(hash, bs)

		err := dps.uncheckedCache.Set(depHash, blocks)
		if err != nil {
			dps.logger.Errorf("Gap previous set cache err for block:%s", hash)
		}
	} else {
		c, err := dps.uncheckedCache.Get(depHash)
		if err != nil {
			dps.logger.Errorf("Gap previous get cache err for block:%s", hash)
		}

		blocks := c.(*sync.Map)
		blocks.Store(hash, bs)
	}
}

func (dps *DPoS) dequeueUnchecked(hash types.Hash) {
	if !dps.uncheckedCache.Has(hash) {
		return
	}

	m, err := dps.uncheckedCache.Get(hash)
	if err != nil {
		dps.logger.Errorf("dequeue unchecked err [%s] for hash [%s]", err, hash)
		return
	}

	cm := m.(*sync.Map)

	cm.Range(func(key, value interface{}) bool {
		bs := value.(*consensus.BlockSource)

		result, _ := dps.lv.Process(bs.Block)
		if result == process.Progress {
			v, e := dps.voteCache.Get(bs.Block.GetHash())
			if e == nil {
				vc := v.(*sync.Map)
				vc.Range(func(key, value interface{}) bool {
					dps.acTrx.vote(value.(*protos.ConfirmAckBlock))
					return true
				})

				dps.voteCache.Remove(bs.Block.GetHash())
			}

			localRepAccount.Range(func(key, value interface{}) bool {
				address := key.(types.Address)
				dps.saveOnlineRep(address)

				va, err := dps.voteGenerate(bs.Block, address, value.(*types.Account))
				if err != nil {
					return true
				}

				dps.acTrx.vote(va)
				dps.eb.Publish(string(common.EventBroadcast), common.ConfirmAck, va)

				return true
			})
		}

		dps.ProcessResult(result, bs)
		return true
	})

	r := dps.uncheckedCache.Remove(hash)
	if !r {
		dps.logger.Error("remove cache for unchecked fail")
	}
}

func (dps *DPoS) ProcessMsg(msgType consensus.MsgType, result process.ProcessResult, bs *consensus.BlockSource, p interface{}) {
	dps.ProcessResult(result, bs)
	hash := bs.Block.GetHash()

	switch msgType {
	case consensus.MsgPublishReq:
		dps.logger.Infof("dps recv publishReq block[%s]", hash)
		if result == process.Progress {
			localRepAccount.Range(func(key, value interface{}) bool {
				address := key.(types.Address)
				dps.saveOnlineRep(address)

				va, err := dps.voteGenerate(bs.Block, address, value.(*types.Account))
				if err != nil {
					return true
				}

				dps.acTrx.vote(va)
				dps.eb.Publish(string(common.EventBroadcast), common.ConfirmAck, va)

				return true
			})
		}
	case consensus.MsgConfirmReq:
		dps.logger.Infof("dps recv confirmReq block[%s]", hash)
		localRepAccount.Range(func(key, value interface{}) bool {
			address := key.(types.Address)
			dps.saveOnlineRep(address)

			va, err := dps.voteGenerate(bs.Block, address, value.(*types.Account))
			if err != nil {
				return true
			}

			if result == process.Progress {
				dps.acTrx.vote(va)
				dps.eb.Publish(string(common.EventBroadcast), common.ConfirmAck, va)
			} else if result == process.Old {
				dps.eb.Publish(string(common.EventBroadcast), common.ConfirmAck, va)
			}

			return true
		})
	case consensus.MsgConfirmAck:
		dps.logger.Infof("dps recv confirmAck block[%s]", hash)
		ack := p.(*protos.ConfirmAckBlock)
		dps.saveOnlineRep(ack.Account)

		//cache the ack messages
		if result == process.GapPrevious || result == process.GapSource {
			if dps.voteCache.Has(hash) {
				v, err := dps.voteCache.Get(hash)
				if err != nil {
					dps.logger.Error("get vote cache err")
					return
				}

				vc := v.(*sync.Map)
				vc.Store(ack.Account, ack)
			} else {
				vc := new(sync.Map)
				vc.Store(ack.Account, ack)
				err := dps.voteCache.Set(hash, vc)
				if err != nil {
					dps.logger.Error("set vote cache err")
					return
				}
			}
		} else if result == process.Progress {
			dps.acTrx.vote(ack)

			localRepAccount.Range(func(key, value interface{}) bool {
				address := key.(types.Address)
				dps.saveOnlineRep(address)

				va, err := dps.voteGenerate(bs.Block, address, value.(*types.Account))
				if err != nil {
					return true
				}

				dps.acTrx.vote(va)
				dps.eb.Publish(string(common.EventBroadcast), common.ConfirmAck, va)

				return true
			})
		} else if result == process.Old {
			dps.acTrx.vote(ack)
		}
	case consensus.MsgSync:
		//
	case consensus.MsgGenerateBlock:
		dps.acTrx.updatePerfTime(hash, time.Now().UnixNano(), false)

		localRepAccount.Range(func(key, value interface{}) bool {
			address := key.(types.Address)
			dps.saveOnlineRep(address)

			va, err := dps.voteGenerate(bs.Block, address, value.(*types.Account))
			if err != nil {
				return true
			}

			dps.acTrx.vote(va)
			dps.eb.Publish(string(common.EventBroadcast), common.ConfirmAck, va)
			return true
		})
	default:
		//
	}
}

func (dps *DPoS) ProcessResult(result process.ProcessResult, bs *consensus.BlockSource) {
	blk := bs.Block
	hash := blk.GetHash()

	switch result {
	case process.Progress:
		if bs.BlockFrom == types.Synchronized {
			dps.logger.Infof("Block %s from sync,no need consensus", hash)
		} else if bs.BlockFrom == types.UnSynchronized {
			dps.logger.Infof("Block %s basic info is correct,begin add it to roots", hash)
			dps.acTrx.addToRoots(blk)
		} else {
			dps.logger.Errorf("Block %s UnKnow from", hash)
		}
		dps.dequeueUnchecked(hash)
	case process.BadSignature:
		dps.logger.Errorf("Bad signature for block: %s", hash)
	case process.BadWork:
		dps.logger.Errorf("Bad work for block: %s", hash)
	case process.BalanceMismatch:
		dps.logger.Errorf("Balance mismatch for block: %s", hash)
	case process.Old:
		dps.logger.Debugf("Old for block: %s", hash)
	case process.UnReceivable:
		dps.logger.Errorf("UnReceivable for block: %s", hash)
	case process.GapSmartContract:
		dps.logger.Errorf("GapSmartContract for block: %s", hash)
		//dps.processGapSmartContract(blk)
	case process.InvalidData:
		dps.logger.Errorf("InvalidData for block: %s", hash)
	case process.Other:
		dps.logger.Errorf("UnKnow process result for block: %s", hash)
	case process.Fork:
		dps.logger.Errorf("Fork for block: %s", hash)
		dps.ProcessFork(blk)
	case process.GapPrevious:
		dps.logger.Debugf("Gap previous for block: %s", hash)
		dps.enqueueUnchecked(hash, blk.Previous, bs)
	case process.GapSource:
		dps.logger.Debugf("Gap source for block: %s", hash)
		dps.enqueueUnchecked(hash, blk.Link, bs)
	}
}

func (dps *DPoS) ProcessFork(block *types.StateBlock) {
	blk := dps.findAnotherForkedBlock(block)
	if _, ok := dps.acTrx.roots.Load(blk.Parent()); !ok {
		dps.acTrx.addToRoots(blk)
		dps.eb.Publish(string(common.EventBroadcast), common.ConfirmReq, blk)
	}
}

func (dps *DPoS) findAnotherForkedBlock(block *types.StateBlock) *types.StateBlock {
	hash := block.Parent()

	forkedHash, err := dps.ledger.GetChild(hash, block.Address)
	if err != nil {
		dps.logger.Error(err)
		return block
	}

	forkedBlock, err := dps.ledger.GetStateBlock(forkedHash)
	if err != nil {
		dps.logger.Error(err)
		return block
	}

	return forkedBlock
}

func (dps *DPoS) voteGenerate(block *types.StateBlock, account types.Address, acc *types.Account) (*protos.ConfirmAckBlock, error) {
	va := &protos.ConfirmAckBlock{
		Sequence:  0,
		Blk:       block,
		Account:   account,
		Signature: acc.Sign(block.GetHash()),
	}
	return va, nil
}

func (dps *DPoS) refreshAccount() {
	var b bool
	var addr types.Address

	for _, v := range dps.accounts {
		addr = v.Address()
		b = dps.isRepresentation(addr)
		if b {
			localRepAccount.Store(addr, v)
		}
	}

	var count uint32
	localRepAccount.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	dps.logger.Infof("there is %d reps", count)
	if count > 1 {
		dps.logger.Error("it is very dangerous to run two or more representatives on one node")
	}
}

func (dps *DPoS) isRepresentation(address types.Address) bool {
	if _, err := dps.ledger.GetRepresentation(address); err != nil {
		return false
	}
	return true
}

func (dps *DPoS) saveOnlineRep(addr types.Address) {
	now := time.Now().Add(repTimeout).UTC().Unix()
	dps.onlineReps.Store(addr, now)
}

func (dps *DPoS) GetOnlineRepresentatives() []types.Address {
	var repAddresses []types.Address

	dps.onlineReps.Range(func(key, value interface{}) bool {
		addr := key.(types.Address)
		repAddresses = append(repAddresses, addr)
		return true
	})

	return repAddresses
}

func (dps *DPoS) findOnlineRepresentatives() error {
	var address types.Address

	localRepAccount.Range(func(key, value interface{}) bool {
		address = key.(types.Address)
		dps.saveOnlineRep(address)
		return true
	})

	blk, err := dps.ledger.GetRandomStateBlock()
	if err != nil {
		return err
	}

	//dps.ns.Broadcast(p2p.ConfirmReq, blk)
	dps.eb.Publish(string(common.EventBroadcast), common.ConfirmReq, blk)
	return nil
}

func (dps *DPoS) cleanOnlineReps() {
	var repAddresses []*types.Address
	now := time.Now().UTC().Unix()

	dps.onlineReps.Range(func(key, value interface{}) bool {
		addr := key.(types.Address)
		v := value.(int64)
		if v < now {
			dps.onlineReps.Delete(addr)
		} else {
			repAddresses = append(repAddresses, &addr)
		}
		return true
	})

	_ = dps.ledger.SetOnlineRepresentations(repAddresses)
}

func (dps *DPoS) calculateAckHash(va *protos.ConfirmAckBlock) (types.Hash, error) {
	data, err := protos.ConfirmAckBlockToProto(va)
	if err != nil {
		return types.ZeroHash, err
	}

	version := dps.cfg.Version
	message := p2p.NewQlcMessage(data, byte(version), common.ConfirmAck)
	hash, err := types.HashBytes(message)
	if err != nil {
		return types.ZeroHash, err
	}

	return hash, nil
}
