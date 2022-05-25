package indexer

import (
	"github.com/angelorc/sinfonia-indexer/chain"
	"github.com/avast/retry-go"
	"log"
	"os"
	"time"
)

var (
	RtyAttNum = uint(5)
	RtyAtt    = retry.Attempts(RtyAttNum)
	RtyDel    = retry.Delay(time.Millisecond * 400)
	RtyErr    = retry.LastErrorOnly(true)
)

func Start() {
	client, err := chain.NewChainClient(chain.GetBitsongConfig(), os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", "bitsong", err)
	}

	indexer := &Indexer{Client: client}
	// run the indexer

	concurrentBlocks := 100
	maxBlock := 10000

	var to_index []int64
	for i := 1; i <= maxBlock; i++ {
		to_index = append(to_index, int64(i))
	}

	if err := indexer.ForEachBlock(to_index, indexer.IndexTransactions, concurrentBlocks); err != nil {
		log.Fatalf("failed to index blocks. err: %v", err)
	}
}
