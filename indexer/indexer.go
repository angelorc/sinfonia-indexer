package indexer

import (
	"context"
	"fmt"
	"github.com/angelorc/sinfonia-indexer/chain"
	"github.com/avast/retry-go"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

type Indexer struct {
	Client *chain.ChainClient
	//DB     *sql.DB
}

func (i *Indexer) ForEachBlock(blocks []int64, cb func(height int64) error, concurrentBlocks int) error {
	fmt.Println("starting block queries for", i.Client.Config.ChainID)

	var (
		eg           errgroup.Group
		mutex        sync.Mutex
		failedBlocks = make([]int64, 0)
		sem          = make(chan struct{}, concurrentBlocks)
	)

	for _, h := range blocks {
		h := h
		sem <- struct{}{}

		eg.Go(func() error {
			if err := cb(h); err != nil {
				if strings.Contains(err.Error(), "wrong ID: no ID") {
					mutex.Lock()
					failedBlocks = append(failedBlocks, h)
					mutex.Unlock()
				} else {
					return fmt.Errorf("[height %d] - failed to get block. err: %s", h, err.Error())
				}
			}
			<-sem
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if len(failedBlocks) > 0 {
		return i.ForEachBlock(failedBlocks, cb, concurrentBlocks)
	}

	return nil
}

func (i *Indexer) IndexTransactions(height int64) error {
	fmt.Println(fmt.Sprintf("Index transactions on block %d", height))
	block, err := i.Client.RPCClient.Block(context.Background(), &height)
	if err != nil {
		if err = retry.Do(func() error {
			block, err = i.Client.RPCClient.Block(context.Background(), &height)
			if err != nil {
				return err
			}

			return nil
		}, RtyAtt, RtyDel, RtyErr, retry.DelayType(retry.BackOffDelay), retry.OnRetry(func(n uint, err error) {
			// indexer.LogRetryGetBlock(n, err, height)
		})); err != nil {
			return err
		}
	}
	if block != nil {
		//i.ParseTxs(block)
		fmt.Println(fmt.Sprintf("Parse txs: %s", block.BlockID))
	}

	return nil
}
