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
        	"token": "a7e8fa30c063e96a489a47bc43909505bd86735da4a109dca28be936118a8582",
        	"address": "qlc_3qjky1ptg9qkzm8iertdzrnx9btjbaea33snh1w4g395xqqczye4kgcfyfs1",
        	"balance": "0",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "bf86c83fb4bfb9f49b9b8fa593c8cf4128c9e21720c487565b52fc6640a9e8f3",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "6TrdxKfo+jDAY+lqSJpHvEOQlQW9hnNdpKEJ3KKL6TYRioWCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADVKa6ehgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAi/hsg/tL+59Jubj6WTyM9BKMniFyDEh1ZbUvxmQKno8wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADUUxDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1FMQwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
        	"povHeight": 0,
        	"timestamp": 1553990401,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_3hw8s1zubhxsykfsq5x7kh6eyibas9j3ga86ixd7pnqwes1cmt9mqqrngap4",
        	"work": "000000000048f5b9",
        	"signature": "18fcd023d9c14bede4a5abad71eca562c8bd5df48ac293466b48b3c2ece42fa00c0e64b68cccf5cfa0f0077ebdc9e05efd757999e3f3c68280975458e9cad40e"
        }`

	testJsonGenesis = `{
        	"type": "ContractReward",
        	"token": "a7e8fa30c063e96a489a47bc43909505bd86735da4a109dca28be936118a8582",
        	"address": "qlc_3hw8s1zubhxsykfsq5x7kh6eyibas9j3ga86ixd7pnqwes1cmt9mqqrngap4",
        	"balance": "60000000000000000",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "8b54787c668dddd4f22ad64a8b0d241810871b9a52a989eb97670f345ad5dc90",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "p+j6MMBj6WpImke8Q5CVBb2Gc12koQncoovpNhGKhYIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANUprp6GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACL+GyD+0v7n0m5uPpZPIz0EoyeIXIMSHVltS/GZAqejzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAL+GyD+0v7n0m5uPpZPIz0EoyeIXIMSHVltS/GZAqejzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANRTEMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADUUxDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        	"povHeight": 0,
        	"timestamp": 1553990410,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_3hw8s1zubhxsykfsq5x7kh6eyibas9j3ga86ixd7pnqwes1cmt9mqqrngap4",
        	"work": "000000000048f5b9",
        	"signature": "a717e690216e357d1b4e200478ef74c51d0b6ab28893fd5cf22aff6f1403d60996ec97bd84214c7a7aed4b0428671e81048afa9b86126c7484b3a88b725e1202"
        }`

	testJsonMintageQGAS = `{
        	"type": "ContractSend",
        	"token": "89066d747a3c74ff1dec8ea6a7011bde010dd404aec454880f23d58cbf9280e4",
        	"address": "qlc_3qjky1ptg9qkzm8iertdzrnx9btjbaea33snh1w4g395xqqczye4kgcfyfs1",
        	"balance": "0",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "e813e51a6d8abea178a2f376d532df983cca71b4e4cf5bdd2d7864ee30cf8ba5",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "6TrdxIkGbXR6PHT/HeyOpqcBG94BDdQErsRUiA8j1Yy/koDkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAjhvJvwQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAjoE+UabYq+oXii83bVMt+YPMpxtOTPW90teGTuMM+LpQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEUUdBUwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABFFHQVMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
        	"povHeight": 0,
        	"timestamp": 1553990401,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_3hw8s1zubhxsykfsq5x7kh6eyibas9j3ga86ixd7pnqwes1cmt9mqqrngap4",
        	"work": "000000000048f5b9",
        	"signature": "441b26cf4318cea394fe07a5e30cde18f967406a9c26158417bcd29abd5a4c79d05746f838bc42f0a7d681cf4a3b4e6b29992fcd7fa7cafe72a4e00e133d310f"
        }`

	testJsonGenesisQGAS = `{
        	"type": "ContractReward",
        	"token": "89066d747a3c74ff1dec8ea6a7011bde010dd404aec454880f23d58cbf9280e4",
        	"address": "qlc_3t1mwnf8u4oyn7wc7wuptnsfz83wsbrubs8hdhgkty56xrrez4x7fcttk5f3",
        	"balance": "10000000000000000",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "f798089896ffdf45ccce2e039666014b8c666ea0f47f0df4ee7e73b49dac0945",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "iQZtdHo8dP8d7I6mpwEb3gEN1ASuxFSIDyPVjL+SgOQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACOG8m/BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACOgT5Rptir6heKLzdtUy35g8ynG05M9b3S14ZO4wz4ulAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOgT5Rptir6heKLzdtUy35g8ynG05M9b3S14ZO4wz4ulAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARRR0FTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEUUdBUwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        	"povHeight": 0,
        	"timestamp": 1553990410,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_3hw8s1zubhxsykfsq5x7kh6eyibas9j3ga86ixd7pnqwes1cmt9mqqrngap4",
        	"work": "000000000048f5b9",
        	"signature": "69903343b5188cedc3b301288a553f3e094bffdf8d1173eb897e630860642a4d62d1022b002c290ca996027d5e424056adad04b340f52d0185362dfd41e07e0c"
        }`

	//test net
	testChainToken, _       = types.NewHash("a7e8fa30c063e96a489a47bc43909505bd86735da4a109dca28be936118a8582")
	testGenesisAddress, _   = types.HexToAddress("qlc_3hw8s1zubhxsykfsq5x7kh6eyibas9j3ga86ixd7pnqwes1cmt9mqqrngap4")
	testGenesisMintageBlock types.StateBlock
	testGenesisMintageHash  types.Hash
	testGenesisBlock        types.StateBlock
	testGenesisBlockHash    types.Hash

	//test net gas
	testGasToken, _     = types.NewHash("89066d747a3c74ff1dec8ea6a7011bde010dd404aec454880f23d58cbf9280e4")
	testGasAddress, _   = types.HexToAddress("qlc_3t1mwnf8u4oyn7wc7wuptnsfz83wsbrubs8hdhgkty56xrrez4x7fcttk5f3")
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