// +build !testnet

/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package config

func DefaultConfigV2(dir string) (*ConfigV2, error) {
	pk, id, err := identityConfig()
	if err != nil {
		return nil, err
	}
	var cfg ConfigV2
	modules := []string{"qlcclassic", "ledger", "account", "net", "util", "wallet", "mintage", "contract", "sms"}
	cfg = ConfigV2{
		Version:             2,
		DataDir:             dir,
		StorageMax:          "10GB",
		AutoGenerateReceive: false,
		LogLevel:            "error",
		PerformanceEnabled:  false,
		RPC: &RPCConfigV2{
			Enable:           false,
			HTTPEnabled:      true,
			HTTPEndpoint:     "tcp4://0.0.0.0:29735",
			HTTPCors:         []string{"*"},
			HttpVirtualHosts: []string{},
			WSEnabled:        true,
			WSEndpoint:       "tcp4://0.0.0.0:29736",
			IPCEnabled:       true,
			IPCEndpoint:      defaultIPCEndpoint(),
			PublicModules:    modules,
		},
		P2P: &P2PConfigV2{
			BootNodes: []string{
				"/ip4/218.17.39.34/tcp/29734/ipfs/QmVYKUm5YPAks18S1K7VenvqsEssWZPh38cSshbdFUhWTC",
			},
			Listen:       "/ip4/0.0.0.0/tcp/29734",
			SyncInterval: 120,
			Discovery: &DiscoveryConfigV2{
				DiscoveryInterval: 10,
				Limit:             20,
				MDNSEnabled:       false,
				MDNSInterval:      30,
			},
			ID: &IdentityConfigV2{id, pk},
		},
	}

	return &cfg, nil
}
