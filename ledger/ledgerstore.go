/*
 * Copyright (c) 2018 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package ledger

import (
	"github.com/qlcchain/go-qlc/common/types"
)

type LedgerStore interface {
	Empty() (bool, error)
	BatchUpdate(fn func() error) error

	// account meta CURD
	AddAccountMeta(meta *types.AccountMeta) error
	GetAccountMeta(address types.Address) (*types.AccountMeta, error)
	UpdateAccountMeta(meta *types.AccountMeta) error
	DeleteAccountMeta(address types.Address) error
	HasAccountMeta(address types.Address) (bool, error)
	// token meta CURD
	AddTokenMeta(address types.Address, meta *types.TokenMeta) error
	GetTokenMeta(address types.Address, tokenType types.Hash) (*types.TokenMeta, error)
	UpdateTokenMeta(address types.Address, meta *types.TokenMeta) error
	DeleteTokenMeta(address types.Address, meta *types.TokenMeta) error
	// blocks CURD
	AddBlock(blk types.Block) error
	GetBlock(hash types.Hash) (types.Block, error)
	GetBlocks() ([]*types.Block, error)
	DeleteBlock(hash types.Hash) error
	HasBlock(hash types.Hash) (bool, error)
	CountBlocks() (uint64, error)
	GetRandomBlock() (types.Block, error)
	// representation CURD
	AddRepresentation(address types.Address, amount types.Balance) error
	SubRepresentation(address types.Address, amount types.Balance) error
	GetRepresentation(address types.Address) (types.Balance, error)
	// unchecked CURD
	AddUncheckedBlock(parentHash types.Hash, blk types.Block, kind types.UncheckedKind) error
	GetUncheckedBlock(parentHash types.Hash, kind types.UncheckedKind) (types.Block, error)
	DeleteUncheckedBlock(parentHash types.Hash, kind types.UncheckedKind) error
	HasUncheckedBlock(hash types.Hash, kind types.UncheckedKind) (bool, error)
	WalkUncheckedBlocks(visit types.UncheckedBlockWalkFunc) error
	CountUncheckedBlocks() (uint64, error)
	// pending CURD
	AddPending(destination types.Address, hash types.Hash, pending *types.PendingInfo) error
	GetPending(destination types.Address, hash types.Hash) (*types.PendingInfo, error)
	DeletePending(destination types.Address, hash types.Hash) error
	// frontier CURD
	AddFrontier(frontier *types.Frontier) error
	GetFrontier(hash types.Hash) (*types.Frontier, error)
	GetFrontiers() ([]*types.Frontier, error)
	DeleteFrontier(hash types.Hash) error
	CountFrontiers() (uint64, error)
}