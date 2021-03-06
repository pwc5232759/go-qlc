/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package contract

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/qlcchain/go-qlc/common"
	"github.com/qlcchain/go-qlc/vm/vmstore"

	"github.com/qlcchain/go-qlc/common/types"
	cabi "github.com/qlcchain/go-qlc/vm/contract/abi"
)

type pledgeInfo struct {
	pledgeTime   *timeSpan
	pledgeAmount *big.Int
}

type Nep5Pledge struct {
}

func (p *Nep5Pledge) GetFee(ctx *vmstore.VMContext, block *types.StateBlock) (types.Balance, error) {
	return types.ZeroBalance, nil
}

// check pledge chain coin
// - address is normal user address
// - small than min pledge amount
// transfer quota to beneficial address
func (*Nep5Pledge) DoSend(ctx *vmstore.VMContext, block *types.StateBlock) error {
	// check pledge amount
	amount, err := ctx.CalculateAmount(block)
	if err != nil {
		return err
	}

	// check send account is user account
	b, err := ctx.IsUserAccount(block.Address)
	if err != nil {
		return err
	}

	if block.Token != common.ChainToken() || amount.IsZero() || !b {
		return errors.New("invalid block data")
	}

	param := new(cabi.PledgeParam)
	if err := cabi.NEP5PledgeABI.UnpackMethod(param, cabi.MethodNEP5Pledge, block.Data); err != nil {
		fmt.Println(err)
		return errors.New("invalid beneficial address")
	}

	pt := cabi.PledgeType(param.PType)
	if info, b := config[pt]; !b {
		return fmt.Errorf("unsupport type %s", pt.String())
	} else if amount.Compare(types.Balance{Int: info.pledgeAmount}) == types.BalanceCompSmaller {
		return fmt.Errorf("not enough pledge amount %s, expect %s", amount.String(), info.pledgeAmount)
	}

	if param.PledgeAddress != block.Address {
		return fmt.Errorf("invalid pledge address[%s],expect %s",
			param.PledgeAddress.String(), block.Address.String())
	}

	block.Data, err = cabi.NEP5PledgeABI.PackMethod(cabi.MethodNEP5Pledge, param.Beneficial,
		param.PledgeAddress, uint8(param.PType), param.NEP5TxId)
	if err != nil {
		return err
	}
	return nil
}

func (*Nep5Pledge) DoReceive(ctx *vmstore.VMContext, block, input *types.StateBlock) ([]*ContractBlock, error) {
	param, err := cabi.ParsePledgeParam(input.Data)
	if err != nil {
		return nil, err
	}
	amount, _ := ctx.CalculateAmount(input)

	var withdrawTime int64
	pt := cabi.PledgeType(param.PType)
	if info, b := config[pt]; b {
		withdrawTime = info.pledgeTime.Calculate(time.Unix(input.Timestamp, 0)).UTC().Unix()
	} else {
		return nil, fmt.Errorf("unsupport type %s", pt.String())
	}

	info := cabi.NEP5PledgeInfo{
		PType:         param.PType,
		Amount:        amount.Int,
		WithdrawTime:  withdrawTime,
		Beneficial:    param.Beneficial,
		PledgeAddress: param.PledgeAddress,
		NEP5TxId:      param.NEP5TxId,
	}

	if _, err := ctx.GetStorage(types.NEP5PledgeAddress[:], []byte(param.NEP5TxId)); err == nil {
		return nil, fmt.Errorf("invalid nep5 tx id")
	} else {
		if err := ctx.SetStorage(types.NEP5PledgeAddress[:], []byte(param.NEP5TxId), nil); err != nil {
			return nil, err
		}
	}

	pledgeKey := cabi.GetPledgeKey(input.Address, param.Beneficial, param.NEP5TxId)

	var pledgeData []byte
	if pledgeData, err = ctx.GetStorage(types.NEP5PledgeAddress[:], pledgeKey); err != nil && err != vmstore.ErrStorageNotFound {
		return nil, err
	} else {
		// already exist,verify data
		if len(pledgeData) > 0 {
			oldPledge, err := cabi.ParsePledgeInfo(pledgeData)
			if err != nil {
				return nil, err
			}
			if oldPledge.PledgeAddress != info.PledgeAddress || oldPledge.WithdrawTime != info.WithdrawTime ||
				oldPledge.Beneficial != info.Beneficial || oldPledge.PType != info.PType ||
				oldPledge.NEP5TxId != info.NEP5TxId {
				return nil, errors.New("invalid saved pledge info")
			}
		} else {
			// save data
			pledgeData, err = cabi.NEP5PledgeABI.PackVariable(cabi.VariableNEP5PledgeInfo, info.PType, info.Amount,
				info.WithdrawTime, info.Beneficial, info.PledgeAddress, info.NEP5TxId)
			if err != nil {
				return nil, err
			}
			err = ctx.SetStorage(types.NEP5PledgeAddress[:], pledgeKey, pledgeData)
			if err != nil {
				return nil, err
			}
		}
	}
	am, _ := ctx.GetAccountMeta(param.Beneficial)
	if am != nil {
		tm := am.Token(common.ChainToken())
		block.Type = types.ContractReward
		block.Address = param.Beneficial
		block.Token = input.Token
		block.Link = input.GetHash()
		block.Data = pledgeData
		block.Balance = am.CoinBalance
		block.Vote = am.CoinVote
		block.Network = am.CoinNetwork
		block.Oracle = am.CoinOracle
		block.Storage = am.CoinStorage
		block.Previous = tm.Header
		block.Representative = tm.Representative
	} else {
		block.Type = types.ContractReward
		block.Address = param.Beneficial
		block.Token = input.Token
		block.Link = input.GetHash()
		block.Data = pledgeData
		block.Vote = types.ZeroBalance
		block.Network = types.ZeroBalance
		block.Oracle = types.ZeroBalance
		block.Storage = types.ZeroBalance
		block.Previous = types.ZeroHash
		block.Representative = input.Representative
	}

	//TODO: query snapshot balance
	switch cabi.PledgeType(param.PType) {
	case cabi.Network:
		block.Network = block.Network.Add(amount)
	case cabi.Oracle:
		block.Oracle = block.Oracle.Add(amount)
	case cabi.Storage:
		block.Storage = block.Storage.Add(amount)
	case cabi.Vote:
		block.Vote = block.Vote.Add(amount)
	default:
		break
	}

	return []*ContractBlock{
		{
			VMContext: ctx,
			Block:     block,
			ToAddress: param.Beneficial,
			BlockType: types.ContractReward,
			Amount:    amount,
			Token:     input.Token,
			Data:      pledgeData,
		},
	}, nil
}

func (*Nep5Pledge) GetRefundData() []byte {
	return []byte{1}
}

type WithdrawNep5Pledge struct {
}

func (*WithdrawNep5Pledge) GetFee(ctx *vmstore.VMContext, block *types.StateBlock) (types.Balance, error) {
	return types.ZeroBalance, nil
}

func (*WithdrawNep5Pledge) DoSend(ctx *vmstore.VMContext, block *types.StateBlock) (err error) {
	if amount, err := ctx.CalculateAmount(block); err != nil {
		return err
	} else {
		if block.Type != types.ContractSend || amount.Compare(types.ZeroBalance) == types.BalanceCompEqual {
			return errors.New("invalid block ")
		}
	}

	param := new(cabi.WithdrawPledgeParam)
	if err := cabi.NEP5PledgeABI.UnpackMethod(param, cabi.MethodWithdrawNEP5Pledge, block.Data); err != nil {
		return errors.New("invalid input data")
	}

	if block.Data, err = cabi.NEP5PledgeABI.PackMethod(cabi.MethodWithdrawNEP5Pledge, param.Beneficial,
		param.Amount, param.PType); err != nil {
		return
	}

	return nil
}

func (*WithdrawNep5Pledge) DoReceive(ctx *vmstore.VMContext, block, input *types.StateBlock) ([]*ContractBlock, error) {
	param := new(cabi.WithdrawPledgeParam)
	err := cabi.NEP5PledgeABI.UnpackMethod(param, cabi.MethodWithdrawNEP5Pledge, input.Data)
	if err != nil {
		return nil, err
	}

	pledgeResults := cabi.SearchBeneficialPledgeInfo(ctx, param)

	if len(pledgeResults) == 0 {
		return nil, errors.New("pledge is not ready")
	}

	//if len(pledgeResults) > 2 {
	//	sort.Slice(pledgeResults, func(i, j int) bool {
	//		return pledgeResults[i].PledgeInfo.WithdrawTime > pledgeResults[j].PledgeInfo.WithdrawTime
	//	})
	//}

	pledgeInfo := pledgeResults[0]

	amount, _ := ctx.CalculateAmount(input)

	var pledgeData []byte
	if pledgeData, err = ctx.GetStorage(nil, pledgeInfo.Key[1:]); err != nil && err != vmstore.ErrStorageNotFound {
		return nil, err
	} else {
		// already exist,verify data
		if len(pledgeData) > 0 {
			oldPledge := new(cabi.NEP5PledgeInfo)
			err := cabi.NEP5PledgeABI.UnpackVariable(oldPledge, cabi.VariableNEP5PledgeInfo, pledgeData)
			if err != nil {
				return nil, err
			}
			if oldPledge.PledgeAddress != pledgeInfo.PledgeInfo.PledgeAddress || oldPledge.WithdrawTime != pledgeInfo.PledgeInfo.WithdrawTime ||
				oldPledge.Beneficial != pledgeInfo.PledgeInfo.Beneficial || oldPledge.PType != pledgeInfo.PledgeInfo.PType ||
				oldPledge.NEP5TxId != pledgeInfo.PledgeInfo.NEP5TxId {
				return nil, errors.New("invalid saved pledge info")
			}

			// TODO: save data or change pledge info state
			err = ctx.SetStorage(nil, pledgeInfo.Key[1:], nil)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("invalid pledge data")
		}
	}

	am, _ := ctx.GetAccountMeta(pledgeInfo.PledgeInfo.PledgeAddress)
	if am == nil {
		return nil, fmt.Errorf("%s do not found", pledgeInfo.PledgeInfo.PledgeAddress.String())
	}
	tm := am.Token(common.ChainToken())
	block.Type = types.ContractReward
	block.Address = pledgeInfo.PledgeInfo.PledgeAddress
	block.Token = input.Token
	block.Link = input.GetHash()
	block.Data = pledgeData
	block.Vote = am.CoinVote
	block.Network = am.CoinNetwork
	block.Oracle = am.CoinOracle
	block.Storage = am.CoinStorage
	block.Previous = tm.Header
	block.Representative = tm.Representative
	block.Balance = am.CoinBalance.Add(amount)

	return []*ContractBlock{
		{
			VMContext: ctx,
			Block:     block,
			ToAddress: pledgeInfo.PledgeInfo.PledgeAddress,
			BlockType: types.ContractReward,
			Amount:    amount,
			Token:     input.Token,
			Data:      pledgeData,
		},
	}, nil
}

func (*WithdrawNep5Pledge) GetRefundData() []byte {
	return []byte{2}
}
