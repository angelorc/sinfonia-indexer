package dataloader

import (
	"github.com/angelorc/sinfonia-indexer/generated/dataloaden"
	"github.com/angelorc/sinfonia-indexer/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// step 1:
//go:generate go run github.com/vektah/dataloaden TransactionLoader string *github.com/angelorc/sinfonia-indexer/model.Transaction

// step 2:
// copy generated file to generated/dataloaden and rename generated package name to dataloaden

func getTransactionLoader() *dataloaden.TransactionLoader {
	maxLimit := 150

	return dataloaden.NewTransactionLoader(
		dataloaden.TransactionLoaderConfig{
			MaxBatch: maxLimit,
			Wait:     1 * time.Millisecond,
			Fetch: func(keys []string) ([]*model.Transaction, []error) {
				var filter model.TransactionWhereDTO

				item := model.Transaction{}

				objectIds := make([]primitive.ObjectID, len(keys))
				for i, k := range keys {
					oid, _ := primitive.ObjectIDFromHex(k)
					objectIds[i] = oid
				}
				customQuery := bson.M{"_id": bson.M{"$in": objectIds}}

				items, err := item.List(&filter, nil, nil, &maxLimit, &customQuery)
				if err != nil {
					return nil, []error{err}
				}

				w := make(map[string]*model.Transaction, len(items))
				for _, item := range items {
					w[item.ID.Hex()] = item
				}

				result := make([]*model.Transaction, len(keys))
				for i, key := range keys {
					result[i] = w[key]
				}

				return result, nil
			},
		},
	)
}
