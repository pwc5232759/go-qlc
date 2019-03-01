/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package contract

import (
	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/ledger"
)

//ContractBlock generated by contract
type ContractBlock struct {
	Block     *types.StateBlock
	ToAddress types.Address
	BlockType types.BlockType
	Amount    types.Balance
	Token     types.Hash
	Data      []byte
}

type ChainContract interface {
	GetFee(ledger *ledger.Ledger, block *types.StateBlock) (types.Balance, error)
	// calc and use quota, check tx data
	DoSend(ledger *ledger.Ledger, block *types.StateBlock) (uint64, error)
	// check status, update state
	DoReceive(ledger *ledger.Ledger, block *types.StateBlock, input *types.StateBlock) ([]*ContractBlock, error)
	// refund data at receive error
	GetRefundData() []byte
	GetQuota() uint64
}