package chain

type ChainClientConfig struct {
	ChainID  string `json:"chain-id" yaml:"chain-id"`
	RPCAddr  string `json:"rpc-addr" yaml:"rpc-addr"`
	GRPCAddr string `json:"grpc-addr" yaml:"grpc-addr"`
	Timeout  string `json:"timeout" yaml:"timeout"`
}

func GetBitsongConfig() *ChainClientConfig {
	return &ChainClientConfig{
		ChainID:  "bitsong-sinfonia-test-1",
		RPCAddr:  "https://rpc.testnet.bitsong.network:443",
		GRPCAddr: "http://142.132.252.143:9090",
		Timeout:  "10s",
	}
}

func GetOsmosisConfig() *ChainClientConfig {
	return &ChainClientConfig{
		ChainID:  "osmosis-sinfonia-test-1",
		RPCAddr:  "https://rpc.osmosis.devnet.bitsong.network:443",
		GRPCAddr: "http://142.132.252.143:10090",
		Timeout:  "10s",
	}
}
