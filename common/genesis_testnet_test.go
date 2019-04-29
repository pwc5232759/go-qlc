// +build  testnet

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

func TestGenesisBlock2(t *testing.T) {
	h, _ := types.NewHash("1e271a57d0988fda4dbd609e4aa1451644d51d7ac98be496de06cd193925964c")

	h2 := testGenesisBlock.GetHash()
	if h2 != h {
		t.Log(util.ToString(testGenesisBlock))
		t.Fatal("invalid genesis block", h2.String(), h.String())
	}

	h3, _ := types.NewHash("30791507160ed430e745a90026de1384361f6117410f4567b3e0844fa415c7cf")
	h4 := testGenesisMintageBlock.GetHash()
	if h3 != h4 {
		t.Log(util.ToIndentString(testGenesisMintageBlock))
		t.Fatal("invalid genesis mintage block", h3.String(), h4.String())
	}
}

func TestGasBlock2(t *testing.T) {
	h, _ := types.NewHash("d9c68cf28f71a1658958459cee34329472246374da7db8ab40623dc5aa61d66f")

	h2 := testGasBlock.GetHash()
	if h2 != h {
		t.Log(util.ToString(testGasBlock))
		t.Fatal("invalid gas block", h2.String(), h.String())
	}

	h3, _ := types.NewHash("9227311c08fd5d11fefbf30844bdd0babce404efdfcd558c7305fad655f513dd")
	h4 := testGasMintageBlock.GetHash()
	if h3 != h4 {
		t.Log(util.ToIndentString(testGasMintageBlock))
		t.Fatal("invalid gas mintage block", h3.String(), h4.String())
	}
}

func TestIsGenesisToken(t *testing.T) {
	h1, _ := types.NewHash("327531148b1a6302632aa7ad6eb369437d8269a08a55b344bd06b514e4e6ae97")
	h2, _ := types.NewHash("3339a985301a9ba7c35e1e15b78f306f9cdb03676436d013218099a9007714e1")
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
