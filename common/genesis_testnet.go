// +build testnet

/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package common

import (
	"encoding/json"

	"github.com/qlcchain/go-qlc/common/types"
)

var (
	testJsonMintage = `{
        	"type": "ContractSend",
        	"token": "3339a985301a9ba7c35e1e15b78f306f9cdb03676436d013218099a9007714e1",
        	"address": "qlc_3qjky1ptg9qkzm8iertdzrnx9btjbaea33snh1w4g395xqqczye4kgcfyfs1",
        	"balance": "0",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "cf4e1dd1d2a0fa27dab8f6c3cbd397f3a02d09d1ee377c68c50147ce71c46131",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "cex33DM5qYUwGpunw14eFbePMG+c2wNnZDbQEyGAmakAdxThAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADVKa6ehgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAjPTh3R0qD6J9q49sPL05fzoC0J0e43fGjFAUfOccRhMQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANRTEMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADUUxDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        	"povHeight": 0,
        	"timestamp": 1553990401,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_3mtg5qax7a9t6zfdjxp5shbshwx17n6x5ujqhjnec1c9ssrwarbjhncxqbrd",
        	"work": "0000000000000000",
        	"signature": "d6e7676c844a75af12d618daacab65905a1ea3b76dd40559c0ea1759039c257381fa844dec36dc7b6a6d116ad21ed4ad682dd41d3177dec6d5c1aae641adb402"
        }`

	testJsonGenesis = `{
        	"type": "ContractReward",
        	"token": "3339a985301a9ba7c35e1e15b78f306f9cdb03676436d013218099a9007714e1",
        	"address": "qlc_3mtg5qax7a9t6zfdjxp5shbshwx17n6x5ujqhjnec1c9ssrwarbjhncxqbrd",
        	"balance": "60000000000000000",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "30791507160ed430e745a90026de1384361f6117410f4567b3e0844fa415c7cf",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "MzmphTAam6fDXh4Vt48wb5zbA2dkNtATIYCZqQB3FOEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANUprp6GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACM9OHdHSoPon2rj2w8vTl/OgLQnR7jd8aMUBR85xxGExAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAM9OHdHSoPon2rj2w8vTl/OgLQnR7jd8aMUBR85xxGExAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1FMQwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANRTEMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
        	"povHeight": 0,
        	"timestamp": 1553990410,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_3mtg5qax7a9t6zfdjxp5shbshwx17n6x5ujqhjnec1c9ssrwarbjhncxqbrd",
        	"work": "0000000000000000",
        	"signature": "cb346b804118ca66b8d95731396d957562d870f32cf24a7a8b0756b595477a4ecc4309846c201582e1633cccddbb74be47b4f7c46f50ecddf90c4de7b2847503"
        }`

	testJsonMintageQGAS = `{
        	"type": "ContractSend",
        	"token": "38a4f0bd6a103401c31d7b1e6e5be098c5380866b1f07bcd8c8d50404a1c2c98",
        	"address": "qlc_3qjky1ptg9qkzm8iertdzrnx9btjbaea33snh1w4g395xqqczye4kgcfyfs1",
        	"balance": "0",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "3410a4152f53fe97cb2491d24fe0a8ec145d28fe5e0e98f7ea0ca17943766d6e",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "cex33Dik8L1qEDQBwx17Hm5b4JjFOAhmsfB7zYyNUEBKHCyYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAjhvJvwQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAg0EKQVL1P+l8skkdJP4KjsFF0o/l4OmPfqDKF5Q3ZtbgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARRR0FTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEUUdBUwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        	"povHeight": 0,
        	"timestamp": 1553990401,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_1f1inickynzykz7kb6gkbzicju1ndnnhwqigm5uyn573h73qeudgarnk6o1z",
        	"work": "0000000000000000",
        	"signature": "eaf9fc7a4b1db168129d018f3be2913f3d35458aeaec8e26b30b3a400db2fe318ded0415769a5605e92c2f8f24cebd8066a4857d44392c258b4c5c717faa970c"
        }`

	testJsonGenesisQGAS = `{
        	"type": "ContractReward",
        	"token": "38a4f0bd6a103401c31d7b1e6e5be098c5380866b1f07bcd8c8d50404a1c2c98",
        	"address": "qlc_1f1inickynzykz7kb6gkbzicju1ndnnhwqigm5uyn573h73qeudgarnk6o1z",
        	"balance": "10000000000000000",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "9227311c08fd5d11fefbf30844bdd0babce404efdfcd558c7305fad655f513dd",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "OKTwvWoQNAHDHXseblvgmMU4CGax8HvNjI1QQEocLJgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACOG8m/BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACDQQpBUvU/6XyySR0k/gqOwUXSj+Xg6Y9+oMoXlDdm1uAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADQQpBUvU/6XyySR0k/gqOwUXSj+Xg6Y9+oMoXlDdm1uAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABFFHQVMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARRR0FTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
        	"povHeight": 0,
        	"timestamp": 1553990410,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_1f1inickynzykz7kb6gkbzicju1ndnnhwqigm5uyn573h73qeudgarnk6o1z",
        	"work": "0000000000000000",
        	"signature": "1a7ddcfe9245165dd5db52821bccd644981a886d9bc4cde5205bacddba9222c937c3eae0b0afb844cca49271e34fc3cc51a4dad80a860b49f44c316d3beacb07"
        }`

	//test net
	testChainToken, _       = types.NewHash("3339a985301a9ba7c35e1e15b78f306f9cdb03676436d013218099a9007714e1")
	testGenesisAddress, _   = types.HexToAddress("qlc_3mtg5qax7a9t6zfdjxp5shbshwx17n6x5ujqhjnec1c9ssrwarbjhncxqbrd")
	testGenesisMintageBlock types.StateBlock
	testGenesisMintageHash  types.Hash
	testGenesisBlock        types.StateBlock
	testGenesisBlockHash    types.Hash

	//test net gas
	testGasToken, _     = types.NewHash("38a4f0bd6a103401c31d7b1e6e5be098c5380866b1f07bcd8c8d50404a1c2c98")
	testGasAddress, _   = types.HexToAddress("qlc_1f1inickynzykz7kb6gkbzicju1ndnnhwqigm5uyn573h73qeudgarnk6o1z")
	testGasMintageBlock types.StateBlock
	testGasMintageHash  types.Hash
	testGasBlock        types.StateBlock
	testGasBlockHash    types.Hash
)

func init() {
	_ = json.Unmarshal([]byte(testJsonMintage), &testGenesisMintageBlock)
	_ = json.Unmarshal([]byte(testJsonGenesis), &testGenesisBlock)
	testGenesisMintageHash = testGenesisMintageBlock.GetHash()
	testGenesisBlockHash = testGenesisBlock.GetHash()
	//test net gas
	_ = json.Unmarshal([]byte(testJsonMintageQGAS), &testGasMintageBlock)
	_ = json.Unmarshal([]byte(testJsonGenesisQGAS), &testGasBlock)
	testGasMintageHash = testGasMintageBlock.GetHash()
	testGasBlockHash = testGasBlock.GetHash()
}

func GenesisAddress() types.Address {
	return testGenesisAddress
}

func GasAddress() types.Address {
	return testGasAddress
}

func ChainToken() types.Hash {
	return testChainToken
}

func GasToken() types.Hash {
	return testGasToken
}

func GenesisMintageBlock() types.StateBlock {
	return testGenesisMintageBlock
}

func GasMintageBlock() types.StateBlock {
	return testGasMintageBlock
}

func GenesisMintageHash() types.Hash {
	return testGenesisMintageHash
}

func GasMintageHash() types.Hash {
	return testGasMintageHash
}

func GenesisBlock() types.StateBlock {
	return testGenesisBlock
}

func GasBlock() types.StateBlock {
	return testGasBlock
}

func GenesisBlockHash() types.Hash {
	return testGenesisBlockHash
}

func GasBlockHash() types.Hash {
	return testGasBlockHash
}

// IsGenesis check block is chain token genesis
func IsGenesisBlock(block *types.StateBlock) bool {
	h := block.GetHash()
	return h == testGenesisMintageHash || h == testGenesisBlockHash || h == testGasMintageHash || h == testGasBlockHash
}

// IsGenesis check token is chain token genesis
func IsGenesisToken(hash types.Hash) bool {
	return hash == testChainToken || hash == testGasToken
}
