package chain

import (
	"github.com/gogo/protobuf/proto"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"io"
	"time"

	appparams "github.com/bitsongofficial/go-bitsong/app/params"
)

type ChainClient struct {
	//log *zap.Logger

	Config    *ChainClientConfig
	RPCClient rpcclient.Client

	Input  io.Reader
	Output io.Writer

	Codec appparams.EncodingConfig
}

func NewChainClient(config *ChainClientConfig, input io.Reader, output io.Writer) (*ChainClient, error) {
	c := &ChainClient{
		//log: log,

		Config: config,
		Input:  input,
		Output: output,
		Codec:  appparams.MakeEncodingConfig(),
	}

	if err := c.Init(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ChainClient) Init() error {
	timeout, _ := time.ParseDuration(c.Config.Timeout)
	rpcClient, err := NewRPCClient(c.Config.RPCAddr, timeout)
	if err != nil {
		return err
	}

	c.RPCClient = rpcClient

	return nil
}

func NewRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = timeout
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}
	return rpcClient, nil
}

func (c *ChainClient) MarshalProto(res proto.Message) ([]byte, error) {
	return c.Codec.Marshaler.MarshalJSON(res)
}
