package consensus

import (
	"github.com/dgraph-io/badger"
	"github.com/qlcchain/go-qlc/common"
	"github.com/qlcchain/go-qlc/common/event"
	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/ledger"
	"github.com/qlcchain/go-qlc/ledger/process"
	"github.com/qlcchain/go-qlc/p2p/protos"
	"time"
)

type Receiver struct {
	eb event.EventBus
	c  *Consensus
}

func NewReceiver(eb event.EventBus) *Receiver {
	r := &Receiver{
		eb: eb,
	}

	return r
}

func (r *Receiver) init(c *Consensus) {
	r.c = c
}

func (r *Receiver) start() error {
	err := r.eb.SubscribeAsync(string(common.EventPublish), r.ReceivePublish, true)
	if err != nil {
		return err
	}

	err = r.eb.SubscribeAsync(string(common.EventConfirmReq), r.ReceiveConfirmReq, true)
	if err != nil {
		return err
	}

	err = r.eb.SubscribeAsync(string(common.EventConfirmAck), r.ReceiveConfirmAck, true)
	if err != nil {
		return err
	}

	err = r.eb.SubscribeAsync(string(common.EventSyncBlock), r.ReceiveSyncBlock, true)
	if err != nil {
		return err
	}

	err = r.eb.SubscribeAsync(string(common.EventGenerateBlock), r.ReceiveGenerateBlock, true)
	if err != nil {
		return err
	}

	return nil
}

func (r *Receiver) stop() error {
	err := r.eb.Unsubscribe(string(common.EventPublish), r.ReceivePublish)
	if err != nil {
		return err
	}

	err = r.eb.Unsubscribe(string(common.EventConfirmReq), r.ReceiveConfirmReq)
	if err != nil {
		return err
	}

	err = r.eb.Unsubscribe(string(common.EventConfirmAck), r.ReceiveConfirmAck)
	if err != nil {
		return err
	}

	err = r.eb.Unsubscribe(string(common.EventSyncBlock), r.ReceiveSyncBlock)
	if err != nil {
		return err
	}

	err = r.eb.Unsubscribe(string(common.EventGenerateBlock), r.ReceiveGenerateBlock)
	if err != nil {
		return err
	}

	return nil
}

func (r *Receiver) ReceivePublish(blk *types.StateBlock, hash types.Hash, msgFrom string) {
	r.c.logger.Infof("receive publish block [%s] from [%s]", blk.GetHash(), msgFrom)
	if !r.c.processed(hash) {
		r.c.processedUpdate(hash)

		var result process.ProcessResult
		var err error

		bs := &BlockSource{
			Block:     blk,
			BlockFrom: types.UnSynchronized,
		}

		result, err = r.Process(bs.Block)
		if err != nil {
			r.c.logger.Errorf("error: [%s] when verify block:[%s]", err, bs.Block.GetHash())
		} else {
			r.eb.Publish(string(common.EventSendMsgToPeers), common.PublishReq, bs.Block, msgFrom)
			r.c.ca.ProcessMsg(MsgPublishReq, result, bs, nil)
		}
	}
}

func (r *Receiver) ReceiveConfirmReq(blk *types.StateBlock, hash types.Hash, msgFrom string) {
	r.c.logger.Infof("receive ConfirmReq block [%s] from [%s]", blk.GetHash(), msgFrom)
	bs := &BlockSource{
		Block:     blk,
		BlockFrom: types.UnSynchronized,
	}

	result, err := r.Process(bs.Block)
	if err != nil {
		r.c.logger.Errorf("error: [%s] when verify block:[%s]", err, bs.Block.GetHash())
	} else {
		r.eb.Publish(string(common.EventSendMsgToPeers), common.ConfirmReq, blk, msgFrom)
		r.c.ca.ProcessMsg(MsgConfirmReq, result, bs, nil)
	}
}

func (r *Receiver) ReceiveConfirmAck(ack *protos.ConfirmAckBlock, hash types.Hash, msgFrom string) {
	r.c.logger.Infof("receive ConfirmAck block [%s] from [%s]", ack.Blk.GetHash(), msgFrom)
	valid := IsAckSignValidate(ack)
	if !valid {
		return
	}

	bs := &BlockSource {
		Block:     ack.Blk,
		BlockFrom: types.UnSynchronized,
	}

	result, err := r.Process(bs.Block)
	if err != nil {
		r.c.logger.Errorf("error: [%s] when verify block:[%s]", err, bs.Block.GetHash())
	} else {
		r.eb.Publish(string(common.EventSendMsgToPeers), common.ConfirmAck, ack, msgFrom)
		r.c.ca.ProcessMsg(MsgConfirmAck, result, bs, ack)
	}
}

func (r *Receiver) ReceiveSyncBlock(blk *types.StateBlock) {
	r.c.logger.Infof("Sync Event for block:[%s]", blk.GetHash())
	bs := &BlockSource{
		Block:     blk,
		BlockFrom: types.Synchronized,
	}

	result, err := r.Process(bs.Block)
	if err != nil {
		r.c.logger.Errorf("error: [%s] when verify block:[%s]", err, bs.Block.GetHash())
	} else {
		r.c.ca.ProcessMsg(MsgSync, result, bs, nil)
	}
}

func (r *Receiver) ReceiveGenerateBlock(result process.ProcessResult, blk *types.StateBlock) {
	r.c.logger.Infof("GenerateBlock Event for block:[%s]", blk.GetHash())
	bs := &BlockSource{
		Block:     blk,
		BlockFrom: types.UnSynchronized,
	}

	r.c.ca.ProcessMsg(MsgGenerateBlock, result, bs, nil)
}

func (r *Receiver) Process(block *types.StateBlock) (process.ProcessResult, error) {
	for {
		if r, err := r.c.verifier.BlockCheck(block); r != process.Progress || err != nil {
			if err == badger.ErrConflict {
				time.Sleep(1 * time.Millisecond)
				continue
			}

			if err != ledger.ErrBlockExists {
				return r, err
			}
		}

		break
	}

	for {
		if err := r.c.verifier.BlockProcess(block); err != nil {
			if err == badger.ErrConflict {
				time.Sleep(1 * time.Millisecond)
				continue
			}

			if err != ledger.ErrBlockExists {
				return process.Other, err
			}
		}

		break
	}

	return process.Progress, nil
}