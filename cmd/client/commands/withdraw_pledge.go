/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/qlcchain/go-qlc/cmd/util"

	"github.com/abiosoft/ishell"
	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/rpc"
	"github.com/qlcchain/go-qlc/rpc/api"
	"github.com/spf13/cobra"
)

func withdrawPledge() {
	var beneficialAccountP string
	var pledgeAccountP string
	var amountP string
	var pTypeP string

	if interactive {
		beneficialAccount := util.Flag{
			Name:  "beneficialAccount",
			Must:  true,
			Usage: "beneficial account private hex string",
		}
		pledgeAccount := util.Flag{
			Name:  "pledgeAccount",
			Must:  true,
			Usage: "pledge account private hex string",
		}
		amount := util.Flag{
			Name:  "amount",
			Must:  true,
			Usage: "pledge amount",
		}
		pType := util.Flag{
			Name:  "pType",
			Must:  true,
			Usage: "pledge type",
		}

		s := &ishell.Cmd{
			Name: "withdrawPledge",
			Help: "withdraw pledge token",
			Func: func(c *ishell.Context) {
				args := []util.Flag{beneficialAccount, pledgeAccount, amount, pType}
				if util.HelpText(c, args) {
					return
				}
				err := util.CheckArgs(c, args)
				if err != nil {
					util.Warn(err)
					return
				}

				beneficialAccountP = util.StringVar(c.Args, beneficialAccount)
				pledgeAccountP = util.StringVar(c.Args, pledgeAccount)
				amountP = util.StringVar(c.Args, amount)
				pTypeP = util.StringVar(c.Args, pType)

				fmt.Println(beneficialAccountP, pledgeAccountP, amountP, pTypeP)
				if err := withdrawPledgeAction(beneficialAccountP, pledgeAccountP, amountP, pTypeP); err != nil {
					util.Warn(err)
					return
				}
			},
		}
		shell.AddCmd(s)
	} else {
		var accountCmd = &cobra.Command{
			Use:   "withdraw pledge",
			Short: "withdraw pledge token",
			Run: func(cmd *cobra.Command, args []string) {
				err := withdrawPledgeAction(beneficialAccountP, pledgeAccountP, amountP, pTypeP)
				if err != nil {
					cmd.Println(err)
				}
			},
		}
		accountCmd.Flags().StringVar(&pledgeAccountP, "pAccount", "", "pledge account private hex string")
		accountCmd.Flags().StringVar(&beneficialAccountP, "bAccount", "", "beneficial account private hex string")
		accountCmd.Flags().StringVar(&amountP, "amount", "", "pledge amount")
		accountCmd.Flags().StringVar(&pTypeP, "pledgeType", "", "pledge type")
		rootCmd.AddCommand(accountCmd)
	}
}

func withdrawPledgeAction(beneficialAccount, pledgeAccount, amount, pType string) error {
	pBytes, err := hex.DecodeString(pledgeAccount)
	if err != nil {
		return err
	}
	p := types.NewAccount(pBytes)

	bBytes, err := hex.DecodeString(beneficialAccount)
	if err != nil {
		return err
	}
	b := types.NewAccount(bBytes)

	client, err := rpc.Dial(endpointP)
	if err != nil {
		return err
	}
	defer client.Close()

	am := types.StringToBalance(amount)

	withdrawPledgeParam := api.WithdrawPledgeParam{
		Beneficial: b.Address(), Amount: am, PType: pType}

	send := types.StateBlock{}
	err = client.Call(&send, "pledge_getWithdrawPledgeBlock", &withdrawPledgeParam)
	if err != nil {
		return err
	}
	sendHash := send.GetHash()
	send.Signature = b.Sign(sendHash)
	var w types.Work
	worker, _ := types.NewWorker(w, send.Root())
	send.Work = worker.NewWork()

	reward := types.StateBlock{}
	err = client.Call(&reward, "pledge_getWithdrawRewardBlock", &send)

	if err != nil {
		return err
	}
	reward.Signature = p.Sign(reward.GetHash())
	var w2 types.Work
	worker2, _ := types.NewWorker(w2, reward.Root())
	reward.Work = worker2.NewWork()

	//TODO: batch process send/reward
	err = client.Call(nil, "ledger_process", &send)
	if err != nil {
		return err
	}

	err = client.Call(nil, "ledger_process", &reward)
	if err != nil {
		return err
	}
	return nil
}
