package p2p

import (
	"errors"
	"time"

	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/ledger"
	"github.com/qlcchain/go-qlc/p2p/protos"
)

//  Message Type
const (
	PublishReq      = "0" //PublishReq
	ConfirmReq      = "1" //ConfirmReq
	ConfirmAck      = "2" //ConfirmAck
	FrontierRequest = "3" //FrontierReq
	FrontierRsp     = "4" //FrontierRsp
	BulkPullRequest = "5" //BulkPullRequest
	BulkPullRsp     = "6" //BulkPullRsp
	BulkPushBlock   = "7" //BulkPushBlock
)

type MessageService struct {
	netService          *QlcService
	quitCh              chan bool
	messageCh           chan Message
	publishMessageCh    chan Message
	confirmReqMessageCh chan Message
	confirmAckMessageCh chan Message
	rspMessageCh        chan Message
	ledger              *ledger.Ledger
	syncService         *ServiceSync
}

// NewService return new Service.
func NewMessageService(netService *QlcService, ledger *ledger.Ledger) *MessageService {
	ms := &MessageService{
		quitCh:              make(chan bool, 6),
		messageCh:           make(chan Message, 65535),
		publishMessageCh:    make(chan Message, 65535),
		confirmReqMessageCh: make(chan Message, 65535),
		confirmAckMessageCh: make(chan Message, 65535),
		rspMessageCh:        make(chan Message, 65535),
		ledger:              ledger,
		netService:          netService,
	}
	ms.syncService = NewSyncService(netService, ledger)
	return ms
}

// Start start message service.
func (ms *MessageService) Start() {
	// register the network handler.
	netService := ms.netService
	netService.Register(NewSubscriber(ms, ms.publishMessageCh, false, PublishReq))
	netService.Register(NewSubscriber(ms, ms.confirmReqMessageCh, false, ConfirmReq))
	netService.Register(NewSubscriber(ms, ms.confirmAckMessageCh, false, ConfirmAck))
	netService.Register(NewSubscriber(ms, ms.messageCh, false, FrontierRequest))
	netService.Register(NewSubscriber(ms, ms.messageCh, false, FrontierRsp))
	netService.Register(NewSubscriber(ms, ms.messageCh, false, BulkPullRequest))
	netService.Register(NewSubscriber(ms, ms.messageCh, false, BulkPullRsp))
	netService.Register(NewSubscriber(ms, ms.messageCh, false, BulkPushBlock))
	// start loop().
	go ms.startLoop()
	go ms.syncService.Start()
	go ms.publishReqLoop()
	go ms.confirmReqLoop()
	go ms.confirmAckLoop()
}

func (ms *MessageService) startLoop() {
	ms.netService.node.logger.Info("Started Message Service.")
	for {
		select {
		case <-ms.quitCh:
			ms.netService.node.logger.Info("Stopped Message Service.")
			return
		case message := <-ms.messageCh:
			switch message.MessageType() {
			case FrontierRequest:
				ms.syncService.onFrontierReq(message)
			case FrontierRsp:
				ms.syncService.checkFrontier(message)
				//ms.syncService.onFrontierRsp(message)
			case BulkPullRequest:
				ms.syncService.onBulkPullRequest(message)
			case BulkPullRsp:
				ms.syncService.onBulkPullRsp(message)
			case BulkPushBlock:
				ms.syncService.onBulkPushBlock(message)
			default:
				ms.netService.node.logger.Error("Received unknown message.")
				time.Sleep(5 * time.Millisecond)
			}
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (ms *MessageService) publishReqLoop() {
	for {
		select {
		case <-ms.quitCh:
			return
		case message := <-ms.publishMessageCh:
			switch message.MessageType() {
			case PublishReq:
				ms.onPublishReq(message)
			default:
				time.Sleep(5 * time.Millisecond)
			}
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (ms *MessageService) confirmReqLoop() {
	for {
		select {
		case <-ms.quitCh:
			return
		case message := <-ms.confirmReqMessageCh:
			switch message.MessageType() {
			case ConfirmReq:
				ms.onConfirmReq(message)
			default:
				time.Sleep(5 * time.Millisecond)
			}
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (ms *MessageService) confirmAckLoop() {
	for {
		select {
		case <-ms.quitCh:
			return
		case message := <-ms.confirmAckMessageCh:
			switch message.MessageType() {
			case ConfirmAck:
				ms.onConfirmAck(message)
			default:
				time.Sleep(5 * time.Millisecond)
			}
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (ms *MessageService) onPublishReq(message Message) {
	if ms.netService.node.cfg.PerformanceEnabled {
		blk, err := protos.PublishBlockFromProto(message.Data())
		if err != nil {
			ms.netService.node.logger.Error(err)
			return
		}
		hash := blk.Blk.GetHash()
		ms.addPerformanceTime(hash)
	}
	ms.netService.node.logger.Infof("receive publish message")
	ms.netService.msgEvent.GetEvent("consensus").Notify(EventPublish, message)
}

func (ms *MessageService) onConfirmReq(message Message) {
	if ms.netService.node.cfg.PerformanceEnabled {
		blk, err := protos.ConfirmReqBlockFromProto(message.Data())
		if err != nil {
			ms.netService.node.logger.Error(err)
			return
		}
		hash := blk.Blk.GetHash()
		ms.addPerformanceTime(hash)
	}

	ms.netService.msgEvent.GetEvent("consensus").Notify(EventConfirmReq, message)
}

func (ms *MessageService) onConfirmAck(message Message) {
	if ms.netService.node.cfg.PerformanceEnabled {
		ack, err := protos.ConfirmAckBlockFromProto(message.Data())
		if err != nil {
			ms.netService.node.logger.Error(err)
			return
		}
		ms.addPerformanceTime(ack.Blk.GetHash())
	}

	ms.netService.msgEvent.GetEvent("consensus").Notify(EventConfirmAck, message)
}

func (ms *MessageService) Stop() {
	//ms.netService.node.logger.Info("stopped message monitor")
	// quit.
	for i := 0; i < 5; i++ {
		ms.quitCh <- true
	}
	ms.syncService.quitCh <- true
	ms.netService.Deregister(NewSubscriber(ms, ms.publishMessageCh, false, PublishReq))
	ms.netService.Deregister(NewSubscriber(ms, ms.confirmReqMessageCh, false, ConfirmReq))
	ms.netService.Deregister(NewSubscriber(ms, ms.confirmAckMessageCh, false, ConfirmAck))
	ms.netService.Deregister(NewSubscriber(ms, ms.messageCh, false, FrontierRequest))
	ms.netService.Deregister(NewSubscriber(ms, ms.messageCh, false, FrontierRsp))
	ms.netService.Deregister(NewSubscriber(ms, ms.messageCh, false, BulkPullRequest))
	ms.netService.Deregister(NewSubscriber(ms, ms.messageCh, false, BulkPullRsp))
	ms.netService.Deregister(NewSubscriber(ms, ms.messageCh, false, BulkPushBlock))
}

func marshalMessage(messageName string, value interface{}) ([]byte, error) {
	switch messageName {
	case PublishReq:
		packet := protos.PublishBlock{
			Blk: value.(*types.StateBlock),
		}
		data, err := protos.PublishBlockToProto(&packet)
		if err != nil {
			return nil, err
		}
		return data, nil
	case ConfirmReq:
		packet := &protos.ConfirmReqBlock{
			Blk: value.(*types.StateBlock),
		}
		data, err := protos.ConfirmReqBlockToProto(packet)
		if err != nil {
			return nil, err
		}
		return data, nil
	case ConfirmAck:
		data, err := protos.ConfirmAckBlockToProto(value.(*protos.ConfirmAckBlock))
		if err != nil {
			return nil, err
		}
		return data, nil
	case FrontierRequest:
		data, err := protos.FrontierReqToProto(value.(*protos.FrontierReq))
		if err != nil {
			return nil, err
		}
		return data, nil
	case FrontierRsp:
		packet := value.(*protos.FrontierResponse)
		data, err := protos.FrontierResponseToProto(packet)
		if err != nil {
			return nil, err
		}
		return data, nil
	case BulkPullRequest:
		data, err := protos.BulkPullReqPacketToProto(value.(*protos.BulkPullReqPacket))
		if err != nil {
			return nil, err
		}
		return data, nil
	case BulkPullRsp:
		PullRsp := &protos.BulkPullRspPacket{
			Blk: value.(*types.StateBlock),
		}
		data, err := protos.BulkPullRspPacketToProto(PullRsp)
		if err != nil {
			return nil, err
		}
		return data, err
	case BulkPushBlock:
		push := &protos.BulkPush{
			Blk: value.(*types.StateBlock),
		}
		data, err := protos.BulkPushBlockToProto(push)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, errors.New("unKnown Message Type")
	}
}

func (ms *MessageService) addPerformanceTime(hash types.Hash) {
	if exit, err := ms.ledger.IsPerformanceTimeExist(hash); !exit && err == nil {
		if b, err := ms.ledger.HasStateBlock(hash); !b && err == nil {
			t := &types.PerformanceTime{
				Hash: hash,
				T0:   time.Now().UnixNano(),
				T1:   0,
				T2:   0,
				T3:   0,
			}
			err = ms.ledger.AddOrUpdatePerformance(t)
			if err != nil {
				ms.netService.node.logger.Error("error when run AddOrUpdatePerformance in onConfirmAck func")
			}
		}
	}
}
