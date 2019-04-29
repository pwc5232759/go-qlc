// +build  !testnet

/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package common

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/common/util"
)

func TestGenesisBlock(t *testing.T) {
	h, _ := types.NewHash("e7bc60ccdd0175d07797b746a383ed9e3185373748fe77a2493636abcb3bf90f")
	genesisBlock := GenesisBlock()
	h2 := genesisBlock.GetHash()
	if h2 != h {
		t.Log(util.ToString(genesisBlock))
		t.Fatal("invalid genesis block", h2.String(), h.String())
	}

	h3, _ := types.NewHash("467cae1b0aba47902945c14aec79f0c5be31d543439a0ed80da9099cb8e4057b")
	h4 := genesisMintageBlock.GetHash()
	if h3 != h4 {
		t.Log(util.ToIndentString(genesisMintageBlock))
		t.Fatal("invalid genesis mintage block", h3.String(), h4.String())
	}
}

func TestGasBlock1(t *testing.T) {
	h, _ := types.NewHash("62be0ce52d9fa4ae07b1fcf69447f952610536ab4697ebf1859d67bfe658b285")
	gasBlock := GasBlock()
	h2 := gasBlock.GetHash()
	if h2 != h {
		t.Log(util.ToString(gasBlock))
		t.Fatal("invalid gas block", h2.String(), h.String())
	}

	h3, _ := types.NewHash("8492fecd1b25467cfbe3eb369b1b65fbc3c206276ab6a3f19e7ef5897a38d393")
	gasMintageBlock := GasMintageBlock()
	h4 := gasMintageBlock.GetHash()
	if h3 != h4 {
		t.Log(util.ToIndentString(gasMintageBlock))
		t.Fatal("invalid gas mintage block", h3.String(), h4.String())
	}
}

func TestIsGenesisToken(t *testing.T) {
	h1, _ := types.NewHash("327531148b1a6302632aa7ad6eb369437d8269a08a55b344bd06b514e4e6ae97")
	h2, _ := types.NewHash("18ceb6779a31caa2948323ef1a57ec59ee5a182939761ae101f5c4e6163efa1a")
	b1 := IsGenesisToken(h1)
	if b1 {
		t.Fatal("h1 should not be Genesis Token")
	}
	b2 := IsGenesisToken(h2)
	if !b2 {
		t.Fatal("h2 should be Genesis Token")
	}
}

func TestBalanceToRaw(t *testing.T) {
	b1 := types.Balance{Int: big.NewInt(2)}
	i, _ := new(big.Int).SetString("200000000", 10)
	b2 := types.Balance{Int: i}

	type args struct {
		b    types.Balance
		unit string
	}
	tests := []struct {
		name    string
		args    args
		want    types.Balance
		wantErr bool
	}{
		{"Mqlc", args{b: b1, unit: "QLC"}, b2, false},
		//{"Mqn1", args{b: b1, unit: "QN1"}, b2, false},
		//{"Mqn3", args{b: b1, unit: "QN3"}, b2, false},
		//{"Mqn5", args{b: b1, unit: "QN5"}, b2, false},
		//{"Mqn6", args{b: b1, unit: "QN6"}, b1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BalanceToRaw(tt.args.b, tt.args.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceToRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BalanceToRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawToBalance(t *testing.T) {
	b1 := types.Balance{Int: big.NewInt(2)}
	i, _ := new(big.Int).SetString("200000000", 10)
	b2 := types.Balance{Int: i}
	type args struct {
		b    types.Balance
		unit string
	}
	tests := []struct {
		name    string
		args    args
		want    types.Balance
		wantErr bool
	}{
		{"Mqlc", args{b: b2, unit: "QLC"}, b1, false},
		//{"Mqn1", args{b: b2, unit: "QN1"}, b1, false},
		//{"Mqn3", args{b: b2, unit: "QN3"}, b1, false},
		//{"Mqn5", args{b: b2, unit: "QN5"}, b1, false},
		//{"Mqn6", args{b: b2, unit: "QN6"}, b2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RawToBalance(tt.args.b, tt.args.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawToBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawToBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}
