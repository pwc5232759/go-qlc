package consensus

import (
	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/p2p/protos"
	"time"
)

const (
	msgCacheSize             = 65536
	msgCacheExpirationTime   = 10 * time.Minute
	blockCacheExpirationTime = 10 * time.Minute
)

type MsgType byte

const (
	MsgPublishReq MsgType = iota
	MsgConfirmReq
	MsgConfirmAck
	MsgSync
	MsgGenerateBlock
)

type BlockSource struct {
	Block     *types.StateBlock
	BlockFrom types.SynchronizedKind
}

func IsAckSignValidate(va *protos.ConfirmAckBlock) bool {
	hash := va.Blk.GetHash()
	verify := va.Account.Verify(hash[:], va.Signature[:])
	return verify
}
