package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/angelorc/sinfonia-indexer/graph/generated"
	"github.com/angelorc/sinfonia-indexer/model"
	"github.com/angelorc/sinfonia-indexer/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *queryResolver) Transaction(ctx context.Context, where *model.TransactionWhereDTO, search *string) (*model.Transaction, error) {
	if where == nil {
		where = &model.TransactionWhereDTO{}
	}
	if search != nil {
		where.OR = utility.MongoSearchFieldParser(model.SEARCH_FILEDS__TRANSACTION, *search)
	}
	item := model.Transaction{}
	item.One(where)
	if item.Hash == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) Transactions(ctx context.Context, where *model.TransactionWhereDTO, search *string, in []*primitive.ObjectID, orderBy *model.TransactionOrderByENUM, skip *int, limit *int) ([]*model.Transaction, error) {
	if where == nil {
		where = &model.TransactionWhereDTO{}
	}
	if search != nil {
		where.OR = utility.MongoSearchFieldParser(model.SEARCH_FILEDS__TRANSACTION, *search)
	}

	// "in" operation for cherrypicking by ids
	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := model.Transaction{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *queryResolver) TransactionCount(ctx context.Context, where *model.TransactionWhereDTO, search *string) (*int, error) {
	t := model.Transaction{}
	if where == nil {
		where = &model.TransactionWhereDTO{}
	}
	count, err := t.Count(where)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
