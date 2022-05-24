package model

import (
	"context"
	"errors"
	"github.com/angelorc/sinfonia-indexer/utility"
	"time"

	"github.com/angelorc/sinfonia-indexer/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * DB Info
 */

const DB_COLLECTION_NAME__TRANSACTION = "Transaction"
const DB_REF_NAME__TRANSACTION = "default"

/**
 * SEARCH regex fields
 */

var SEARCH_FILEDS__TRANSACTION = []string{"hash"}

/**
 * MODEL
 */

type Transaction struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Hash      string             `json:"hash" bson:"hash" validate:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

/**
 * ENUM
 */

type TransactionOrderByENUM string

/**
 * DTO
 */

// Read

type TransactionWhereUniqueDTO struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type TransactionWhereDTO struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Hash      *string             `json:"hash" bson:"hash,omitempty"`
	CreatedAt time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	OR        []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type UserCreateDTO struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Hash      *string             `json:"hash" bson:"hash,omitempty"`
	CreatedAt time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UserUpdateDTO struct {
	Hash      *string   `json:"hash" bson:"hash,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

/**
 * OPERATIONS
 */

// Read

func (t *Transaction) One(filter *TransactionWhereDTO) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&t)

	return nil
}

func (t *Transaction) List(filter *TransactionWhereDTO, orderBy *TransactionOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Transaction, error) {
	var items []*Transaction
	orderByKey := "created_at"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := options.Find()
	if limit != nil {
		options.SetLimit(int64(*limit))
	}
	if skip != nil {
		options.SetSkip(int64(*skip))
	}
	if orderBy != nil {
		orderByKey, orderByValue = utility.GetOrderByKeyAndValue(string(*orderBy))
	}
	options.SetSort(map[string]int{orderByKey: orderByValue})

	var queryFilter interface{}
	if filter != nil {
		queryFilter = filter
	}
	if !utility.IsZeroVal(customQuery) {
		queryFilter = customQuery
	}

	cursor, err := collection.Find(ctx, &queryFilter, options)
	if err != nil {
		return items, err
	}
	err = cursor.All(ctx, &items)
	if err != nil {
		return items, err
	}

	return items, nil
}

func (t *Transaction) Count(filter *TransactionWhereDTO) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (t *Transaction) Create(data *TransactionWhereDTO) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check uniqe
	item := new(Transaction)
	f := bson.M{
		"$or": []bson.M{
			{"hash": data.Hash},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.Hash != "" {
		return errors.New("user is already exist")
	}
	// operation
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("server error")
	}
	collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&t)
	return nil
}

func (t *Transaction) Update(where primitive.ObjectID, data *UserUpdateDTO) error {
	data.UpdatedAt = time.Now()

	// validate
	if utility.IsZeroVal(where) {
		return errors.New("internal server error")
	}
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check transaction is exists
	collection.FindOne(ctx, bson.M{"_id": where}).Decode(&t)
	if t.Hash == "" {
		return errors.New("item not found")
	}

	// check unique
	item := new(Transaction)
	f := bson.M{
		"$or": []bson.M{
			{"hash": data.Hash, "_id": bson.M{"$ne": where}},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.Hash != "" {
		return errors.New("transaction already exist")
	}

	// operation
	_, err := collection.UpdateOne(ctx, bson.M{"_id": where}, bson.M{"$set": data})
	collection.FindOne(ctx, bson.M{"_id": where}).Decode(&t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) Delete() error {
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if utility.IsZeroVal(t.ID) {
		return errors.New("invalid id")
	}

	collection.FindOne(ctx, bson.M{"_id": t.ID}).Decode(&t)
	if t.Hash == "" {
		return errors.New("item not found")
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": t.ID})
	if err != nil {
		return err
	}

	return nil
}
