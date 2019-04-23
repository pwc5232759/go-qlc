package dpos

import (
	"sync"
	"time"

	"github.com/qlcchain/go-qlc/common"

	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/p2p/protos"
)

const (
	//announcementMin        			= 4 //Minimum number of block announcements
	announcementMax  = 20 //Max number of block announcements
	announceInterval = 16
)

type ActiveTrx struct {
	confirmed electionStatus
	dps       *DPoS
	roots     *sync.Map
	inactive  []types.Hash
	quitCh    chan bool
}

func NewActiveTrx() *ActiveTrx {
	return &ActiveTrx{
		roots: new(sync.Map),
	}
}

func (act *ActiveTrx) SetDposService(dps *DPoS) {
	act.dps = dps
}

func (act *ActiveTrx) start() {
	go act.announceVotes()
}

func (act *ActiveTrx) stop() {
	act.quitCh <- true
}

func (act *ActiveTrx) getVoteKey(block *types.StateBlock) []byte {
	var key [1 + types.HashSize]byte

	if block.Type.Equal(types.Open) || block.Type.Equal(types.ContractReward) {
		key[0] = 1
		copy(key[1:], block.Link[:])
	} else {
		key[0] = 0
		copy(key[1:], block.Previous[:])
	}

	return key[:]
}

func (act *ActiveTrx) addToRoots(block *types.StateBlock) bool {
	vk := act.getVoteKey(block)

	if _, ok := act.roots.Load(vk); !ok {
		ele, err := NewElection(act.dps, block)
		if err != nil {
			act.dps.logger.Errorf("block :%s add to roots error", block.GetHash())
			return false
		}
		act.roots.Store(vk, ele)
		return true
	} else {
		act.dps.logger.Infof("block :%s already exit in roots", block.GetHash())
		return false
	}
}

func (act *ActiveTrx) updatePerformanceTime(hash types.Hash, el *Election) {
	if !act.dps.cfg.PerformanceEnabled {
		return
	}

	if el.announcements == 0 {
		if p, err := act.dps.ledger.GetPerformanceTime(hash); p != nil && err == nil {
			t := &types.PerformanceTime{
				Hash: hash,
				T0:   p.T0,
				T1:   p.T1,
				T2:   time.Now().UnixNano(),
				T3:   p.T3,
			}

			act.dps.ledger.AddOrUpdatePerformance(t)
			if err != nil {
				act.dps.logger.Info("AddOrUpdatePerformance error T2")
			}
		} else {
			act.dps.logger.Info("get performanceTime error T2")
		}
	}

	if el.confirmed {
		var t *types.PerformanceTime
		if p, err := act.dps.ledger.GetPerformanceTime(hash); p != nil && err == nil {
			if el.announcements == 0 {
				t = &types.PerformanceTime{
					Hash: hash,
					T0:   p.T0,
					T1:   time.Now().UnixNano(),
					T2:   p.T2,
					T3:   time.Now().UnixNano(),
				}
			} else {
				t = &types.PerformanceTime{
					Hash: hash,
					T0:   p.T0,
					T1:   time.Now().UnixNano(),
					T2:   p.T2,
					T3:   p.T3,
				}
			}
			err := act.dps.ledger.AddOrUpdatePerformance(t)
			if err != nil {
				act.dps.logger.Info("AddOrUpdatePerformance error T1")
			}
		} else {
			act.dps.logger.Info("get performanceTime error T1")
		}
	} else {
		if el.announcements == 0 {
			if p, err := act.dps.ledger.GetPerformanceTime(hash); p != nil && err == nil {
				t := &types.PerformanceTime{
					Hash: hash,
					T0:   p.T0,
					T1:   p.T1,
					T2:   p.T2,
					T3:   time.Now().UnixNano(),
				}
				act.dps.ledger.AddOrUpdatePerformance(t)
				if err != nil {
					act.dps.logger.Info("AddOrUpdatePerformance error T3")
				}
			} else {
				act.dps.logger.Info("get performanceTime error T3")
			}
		}
	}
}

func (act *ActiveTrx) announceVotes() {
	for {
		nowTime := time.Now().Unix()

		act.roots.Range(func(key, value interface{}) bool {
			el := value.(*Election)
			if nowTime - el.lastTime < announceInterval {
				return true
			} else {
				el.lastTime = nowTime
			}

			block := el.status.winner
			hash := block.GetHash()

			act.updatePerformanceTime(hash, el)

			if el.confirmed {
				act.dps.logger.Infof("block [%s] is already confirmed", hash)
				act.dps.eb.Publish(string(common.EventConfirmedBlock), block)
				act.inactive = append(act.inactive, el.vote.id)
				act.rollBack(el.status.loser)
				act.addWinner2Ledger(block)
			} else {
				act.dps.logger.Infof("vote:send confirmReq for block [%s]", hash)
				act.dps.eb.Publish(string(common.EventBroadcast), common.ConfirmReq, block)
				el.announcements++
			}

			if el.announcements == announcementMax {
				if _, ok := act.roots.Load(value); !ok {
					act.inactive = append(act.inactive, el.vote.id)
				}
			}

			return true
		})

		for _, value := range act.inactive {
			if _, ok := act.roots.Load(value); ok {
				act.roots.Delete(value)
			}
		}
		act.inactive = act.inactive[:0:0]

		time.Sleep(1 * time.Second)
	}
}

func (act *ActiveTrx) addWinner2Ledger(block *types.StateBlock) {
	hash := block.GetHash()

	if exist, err := act.dps.ledger.HasStateBlock(hash); !exist && err == nil {
		err := act.dps.lv.BlockProcess(block)
		if err != nil {
			act.dps.logger.Error(err)
		} else {
			act.dps.logger.Debugf("save block[%s]", hash.String())
		}
	} else {
		act.dps.logger.Debugf("%s, %v", hash.String(), err)
	}
}

func (act *ActiveTrx) rollBack(blocks []*types.StateBlock) {
	for _, v := range blocks {
		hash := v.GetHash()
		act.dps.logger.Info("loser hash is :", hash.String())
		h, err := act.dps.ledger.HasStateBlock(hash)
		if err != nil {
			act.dps.logger.Errorf("error [%s] when run HasStateBlock func ", err)
			continue
		}
		if h {
			err = act.dps.ledger.Rollback(hash)
			if err != nil {
				act.dps.logger.Errorf("error [%s] when rollback hash [%s]", err, hash.String())
			}
		}
	}
}

func (act *ActiveTrx) vote(va *protos.ConfirmAckBlock) {
	vk := act.getVoteKey(va.Blk)

	if v, ok := act.roots.Load(vk); ok {
		v.(*Election).voteAction(va)
	}
}
