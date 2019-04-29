// +build  !testnet

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
	jsonMintage = `{
        	"type": "ContractSend",
        	"token": "18ceb6779a31caa2948323ef1a57ec59ee5a182939761ae101f5c4e6163efa1a",
        	"address": "qlc_3qjky1ptg9qkzm8iertdzrnx9btjbaea33snh1w4g395xqqczye4kgcfyfs1",
        	"balance": "0",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "8699c7c72bd62f24456d0449e59cb735140058fdaf95c9f90055397eb6fd401f",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "cex33BjOtneaMcqilIMj7xpX7FnuWhgpOXYa4QH1xOYWPvoaAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADVKa6ehgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAiGmcfHK9YvJEVtBEnlnLc1FABY/a+VyfkAVTl+tv1AHwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANRTEMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADUUxDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        	"povHeight": 0,
        	"timestamp": 1553990401,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_33nsrz5kqojh6j4pt34bwpgdgfan13ehudwos9wi1obshtuhti1z7wodfb7p",
        	"work": "0000000000000000",
        	"signature": "4e0275c50923bf268e81b5c5eadf9b82c134fbe7f1e2c5443558dc4d740b215ea6c62b71c94951c7d9d1f657073477fa69461b7ff43349b35d1f712566720c03"
        }`
	jsonGenesis = `{
        	"type": "ContractReward",
        	"token": "18ceb6779a31caa2948323ef1a57ec59ee5a182939761ae101f5c4e6163efa1a",
        	"address": "qlc_33nsrz5kqojh6j4pt34bwpgdgfan13ehudwos9wi1obshtuhti1z7wodfb7p",
        	"balance": "60000000000000000",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "467cae1b0aba47902945c14aec79f0c5be31d543439a0ed80da9099cb8e4057b",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "GM62d5oxyqKUgyPvGlfsWe5aGCk5dhrhAfXE5hY++hoAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANUprp6GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACIaZx8cr1i8kRW0ESeWctzUUAFj9r5XJ+QBVOX62/UAfAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIaZx8cr1i8kRW0ESeWctzUUAFj9r5XJ+QBVOX62/UAfAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1FMQwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANRTEMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
        	"povHeight": 0,
        	"timestamp": 1553990410,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_33nsrz5kqojh6j4pt34bwpgdgfan13ehudwos9wi1obshtuhti1z7wodfb7p",
        	"work": "0000000000000000",
        	"signature": "a06e1ec371aed0ebfd467c936af31f07b243c660d9c2b1970bd7bba409cfe28dca8beb423a96c9d4d7d1a3dadc6b0046ba9b3eb91257b0fa1d0a00e95dd96b00"
        }`

	jsonMintageQGAS = `{
        	"type": "ContractSend",
        	"token": "66e27e8c0e82031e3bf6b89a5be4902bcf981725d802b1a3e31fc1c32358b331",
        	"address": "qlc_3qjky1ptg9qkzm8iertdzrnx9btjbaea33snh1w4g395xqqczye4kgcfyfs1",
        	"balance": "0",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "106ccdbff8be60b3b240ba9e5f78b26d59e2ac3bba38149d7a5b93103dcc735d",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "cex33GbifowOggMeO/a4mlvkkCvPmBcl2AKxo+MfwcMjWLMxAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAjhvJvwQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgQbM2/+L5gs7JAup5feLJtWeKsO7o4FJ16W5MQPcxzXQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARRR0FTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEUUdBUwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        	"povHeight": 0,
        	"timestamp": 1553990401,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_165espzzjhm1pgs63gnydxwd6ucswcp5qgjr4kgqnpwm41ywrwtxmd9pyrw4",
        	"work": "0000000000000000",
        	"signature": "467ba381f7278f0b81d3283ef3a9f3b7fc31d595de50f6ffea78d267a81efe5d9870623cb51ecc99d48ff55b0b5a33264a58f1892ab7fb9fa272c692ad6c5e01"
        }`
	jsonGenesisQGAS = `{
        	"type": "ContractReward",
        	"token": "66e27e8c0e82031e3bf6b89a5be4902bcf981725d802b1a3e31fc1c32358b331",
        	"address": "qlc_165espzzjhm1pgs63gnydxwd6ucswcp5qgjr4kgqnpwm41ywrwtxmd9pyrw4",
        	"balance": "10000000000000000",
        	"vote": "0",
        	"network": "0",
        	"storage": "0",
        	"oracle": "0",
        	"previous": "0000000000000000000000000000000000000000000000000000000000000000",
        	"link": "8492fecd1b25467cfbe3eb369b1b65fbc3c206276ab6a3f19e7ef5897a38d393",
        	"message": "0000000000000000000000000000000000000000000000000000000000000000",
        	"data": "ZuJ+jA6CAx479riaW+SQK8+YFyXYArGj4x/BwyNYszEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACOG8m/BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACBBszb/4vmCzskC6nl94sm1Z4qw7ujgUnXpbkxA9zHNdAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBszb/4vmCzskC6nl94sm1Z4qw7ujgUnXpbkxA9zHNdAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABFFHQVMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARRR0FTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
        	"povHeight": 0,
        	"timestamp": 1553990410,
        	"extra": "0000000000000000000000000000000000000000000000000000000000000000",
        	"representative": "qlc_165espzzjhm1pgs63gnydxwd6ucswcp5qgjr4kgqnpwm41ywrwtxmd9pyrw4",
        	"work": "0000000000000000",
        	"signature": "2a49f3542034b98d09b23e28297a77846ac9b6f88fb4e474e914f191d05b38f36088ede35baec6b96e7085944f73e0be1fe2a36acbe00e648edebe0e62c87704"
        }`

	//main net
	chainToken, _       = types.NewHash("18ceb6779a31caa2948323ef1a57ec59ee5a182939761ae101f5c4e6163efa1a")
	genesisAddress, _   = types.HexToAddress("qlc_33nsrz5kqojh6j4pt34bwpgdgfan13ehudwos9wi1obshtuhti1z7wodfb7p")
	genesisMintageBlock types.StateBlock
	genesisMintageHash  types.Hash
	genesisBlock        types.StateBlock
	genesisBlockHash    types.Hash

	//main net gas
	gasToken, _     = types.NewHash("66e27e8c0e82031e3bf6b89a5be4902bcf981725d802b1a3e31fc1c32358b331")
	gasAddress, _   = types.HexToAddress("qlc_165espzzjhm1pgs63gnydxwd6ucswcp5qgjr4kgqnpwm41ywrwtxmd9pyrw4")
	gasMintageBlock types.StateBlock
	gasMintageHash  types.Hash
	gasBlock        types.StateBlock
	gasBlockHash    types.Hash
)

func init() {
	_ = json.Unmarshal([]byte(jsonMintage), &genesisMintageBlock)
	_ = json.Unmarshal([]byte(jsonGenesis), &genesisBlock)
	genesisMintageHash = genesisMintageBlock.GetHash()
	genesisBlockHash = genesisBlock.GetHash()
	//main net gas
	_ = json.Unmarshal([]byte(jsonMintageQGAS), &gasMintageBlock)
	_ = json.Unmarshal([]byte(jsonGenesisQGAS), &gasBlock)
	gasMintageHash = gasMintageBlock.GetHash()
	gasBlockHash = gasBlock.GetHash()
}

func GenesisAddress() types.Address {
	return genesisAddress
}

func GasAddress() types.Address {
	return gasAddress

}

func ChainToken() types.Hash {
	return chainToken
}

func GasToken() types.Hash {
	return gasToken
}

func GenesisMintageBlock() types.StateBlock {
	return genesisMintageBlock
}

func GasMintageBlock() types.StateBlock {
	return gasMintageBlock
}

func GenesisMintageHash() types.Hash {
	return genesisMintageHash
}

func GasMintageHash() types.Hash {
	return gasMintageHash
}

func GenesisBlock() types.StateBlock {
	return genesisBlock
}

func GasBlock() types.StateBlock {
	return gasBlock

}

func GenesisBlockHash() types.Hash {
	return genesisBlockHash
}

func GasBlockHash() types.Hash {
	return gasBlockHash
}

// IsGenesis check block is chain token genesis
func IsGenesisBlock(block *types.StateBlock) bool {
	h := block.GetHash()
	return h == genesisMintageHash || h == genesisBlockHash || h == gasMintageHash || h == gasBlockHash
}

// IsGenesis check token is chain token genesis
func IsGenesisToken(hash types.Hash) bool {
	return hash == chainToken || hash == gasToken
}
