package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// InsertOperation : interface, defines the methods to insert one or many documents in the database
type InsertOperation interface {
	Context(ctx context.Context) InsertOperation
	StopOnFail(stop bool) InsertOperation
	One(document interface{}) (interface{}, error)
	Many(documents []interface{}) ([]interface{}, error)
}

// MongoInsertOperation : implements InsertOperation for MongoDB, it is used to insert one or many documents in the database
type MongoInsertOperation struct {
	collection    *mongo.Collection
	ctx           context.Context
	optInsertMany *options.InsertManyOptions
	err           error
}

// Context : sets the context for the operation
func (fo MongoInsertOperation) Context(ctx context.Context) InsertOperation {
	fo.ctx = ctx
	return fo
}

// StopOnFail : sets the option to stop on fail
func (fo MongoInsertOperation) StopOnFail(stop bool) InsertOperation {
	fo.optInsertMany.SetOrdered(stop)
	return fo
}

// One : inserts one document in the database
func (fo MongoInsertOperation) One(document interface{}) (interface{}, error) {
	if fo.err != nil {
		return nil, fo.err
	}

	result, err := fo.collection.InsertOne(fo.ctx, document)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// Many : inserts many documents in the database
func (fo MongoInsertOperation) Many(documents []interface{}) ([]interface{}, error) {
	if fo.err != nil {
		return nil, fo.err
	}

	results, err := fo.collection.InsertMany(fo.ctx, documents, fo.optInsertMany)
	if err != nil {
		return nil, err
	}

	return results.InsertedIDs, nil
}

// Insert : returns an InsertOperation instance
func (storage *MongoStorage) Insert(collectionName string) InsertOperation {
	collection, ok := storage.collections[collectionName]
	op := MongoInsertOperation{
		collection:    collection,
		optInsertMany: options.InsertMany(),
	}
	if !ok {
		op.err = errors.New("invalid collection: " + collectionName)
	}
	return op
}

// UpdateOperation : interface, defines the methods to update one or many documents in the database
type UpdateOperation interface {
	Context(ctx context.Context) UpdateOperation
	Filter(filter OrderedDocument) UpdateOperation
	Upsert(upsertResults *interface{}) UpdateOperation
	One(update interface{}) (bool, error)
	OneComplex(update interface{}) (bool, error)
	Many(update interface{}) (int64, error)
}

// MongoUpdateOperation : implements UpdateOperation for MongoDB, it is used to update one or many documents in the database
type MongoUpdateOperation struct {
	collection   *mongo.Collection
	filter       OrderedDocument
	update       OrderedDocument
	ctx          context.Context
	opt          *options.UpdateOptions
	upsertResult *interface{}
	err          error
}

// Context : sets the context for the operation
func (fo MongoUpdateOperation) Context(ctx context.Context) UpdateOperation {
	fo.ctx = ctx
	return fo
}

// Filter : sets the filter for the operation
func (fo MongoUpdateOperation) Filter(filter OrderedDocument) UpdateOperation {
	for _, elem := range filter {
		fo.filter = append(fo.filter, primitive.E{Key: elem.Key, Value: elem.Value})
	}
	return fo
}

// Upsert : sets the option to upsert
func (fo MongoUpdateOperation) Upsert(upsertResults *interface{}) UpdateOperation {
	fo.upsertResult = upsertResults
	fo.opt.SetUpsert(true)
	return fo
}

// One : updates one document in the database
func (fo MongoUpdateOperation) One(update interface{}) (bool, error) {
	if fo.err != nil {
		return false, fo.err
	}

	for i := range fo.update {
		fo.update[i].Value = update
	}
	result, err := fo.collection.UpdateOne(fo.ctx, fo.filter, fo.update, fo.opt)
	if err != nil {
		return false, err
	}

	if fo.upsertResult != nil {
		*(fo.upsertResult) = result.UpsertedID
	}
	return result.ModifiedCount == 1, nil
}

// OneComplex : updates one document in the database with the option to upsert
func (fo MongoUpdateOperation) OneComplex(update interface{}) (bool, error) {
	if fo.err != nil {
		return false, fo.err
	}

	result, err := fo.collection.UpdateOne(fo.ctx, fo.filter, update, fo.opt)
	if err != nil {
		return false, err
	}

	if fo.upsertResult != nil {
		*(fo.upsertResult) = result.UpsertedID
	}
	return result.ModifiedCount == 1, nil
}

// Many : updates many documents in the database
func (fo MongoUpdateOperation) Many(update interface{}) (int64, error) {
	if fo.err != nil {
		return 0, fo.err
	}

	for i := range fo.update {
		fo.update[i].Value = update
	}
	result, err := fo.collection.UpdateMany(fo.ctx, fo.filter, fo.update, fo.opt)
	if err != nil {
		return 0, err
	}

	if fo.upsertResult != nil {
		*(fo.upsertResult) = result.UpsertedID
	}
	return result.ModifiedCount, nil
}

// Update : returns an UpdateOperation instance
func (storage *MongoStorage) Update(collectionName string) UpdateOperation {
	collection, ok := storage.collections[collectionName]
	op := MongoUpdateOperation{
		collection: collection,
		filter:     OrderedDocument{},
		update:     OrderedDocument{{"$set", nil}},
		opt:        options.Update(),
	}
	if !ok {
		op.err = errors.New("invalid collection: " + collectionName)
	}
	return op
}

// FindOperation : interface, defines the methods to find one or many documents in the database
type FindOperation interface {
	Context(ctx context.Context) FindOperation
	Filter(filter OrderedDocument) FindOperation
	Projection(filter OrderedDocument) FindOperation
	Sort(field string, ascending bool) FindOperation
	Limit(n int64) FindOperation
	Skip(n int64) FindOperation
	MaxTime(duration time.Duration) FindOperation
	First(result interface{}) error
	All(results interface{}) error
}

// MongoFindOperation : implements FindOperation for MongoDB, it is used to find one or many documents in the database
type MongoFindOperation struct {
	collection *mongo.Collection
	filter     OrderedDocument
	projection OrderedDocument
	ctx        context.Context
	optFind    *options.FindOptions
	optFindOne *options.FindOneOptions
	sorts      []Entry
	err        error
}

// Context : sets the context for the operation
func (fo MongoFindOperation) Context(ctx context.Context) FindOperation {
	fo.ctx = ctx
	return fo
}

// Filter : sets the filter for the operation
func (fo MongoFindOperation) Filter(filter OrderedDocument) FindOperation {
	for _, elem := range filter {
		fo.filter = append(fo.filter, primitive.E{Key: elem.Key, Value: elem.Value})
	}
	return fo
}

// Projection : sets the projection for the operation
func (fo MongoFindOperation) Projection(projection OrderedDocument) FindOperation {
	for _, elem := range projection {
		fo.projection = append(fo.projection, primitive.E{Key: elem.Key, Value: elem.Value})
	}
	fo.optFindOne.SetProjection(fo.projection)
	fo.optFind.SetProjection(fo.projection)
	return fo
}

// Limit : sets the limit for the operation
// limits the results extracted from the database to n
func (fo MongoFindOperation) Limit(n int64) FindOperation {
	fo.optFind.SetLimit(n)
	return fo
}

// Skip : sets the skip option for the operation
// skips the first n results extracted from the database
func (fo MongoFindOperation) Skip(n int64) FindOperation {
	fo.optFind.SetSkip(n)
	return fo
}

// MaxTime : sets the max time for the operation
func (fo MongoFindOperation) MaxTime(duration time.Duration) FindOperation {
	fo.optFind.SetMaxTime(duration)
	return fo
}

// Sort : sets the sort option for the operation
// sorts the results extracted from the database by the specified field in ascending or descending order
func (fo MongoFindOperation) Sort(field string, ascending bool) FindOperation {
	var sort int
	if ascending {
		sort = 1
	} else {
		sort = -1
	}
	fo.sorts = append(fo.sorts, primitive.E{Key: field, Value: sort})
	fo.optFind.SetSort(fo.sorts)
	fo.optFindOne.SetSort(fo.sorts)
	return fo
}

// First : finds the first document in the database
func (fo MongoFindOperation) First(result interface{}) error {
	if fo.err != nil {
		return fo.err
	}

	err := fo.collection.FindOne(fo.ctx, fo.filter, fo.optFindOne).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			result = nil
			return nil
		}

		return err
	}
	return nil
}

// All : finds all the documents in the database
func (fo MongoFindOperation) All(results interface{}) error {
	if fo.err != nil {
		return fo.err
	}
	cursor, err := fo.collection.Find(fo.ctx, fo.filter, fo.optFind)
	if err != nil {
		return err
	}
	err = cursor.All(fo.ctx, results)
	if err != nil {
		return err
	}
	return nil
}

// Find : returns a FindOperation instance
func (storage *MongoStorage) Find(collectionName string) FindOperation {
	collection, ok := storage.collections[collectionName]
	op := MongoFindOperation{
		collection: collection,
		filter:     OrderedDocument{},
		projection: OrderedDocument{},
		optFind:    options.Find(),
		optFindOne: options.FindOne(),
		sorts:      OrderedDocument{},
	}
	if !ok {
		op.err = errors.New("invalid collection: " + collectionName)
	}
	return op
}

// DeleteOperation : interface, defines the methods to delete one or many documents in the database
type DeleteOperation interface {
	Context(ctx context.Context) DeleteOperation
	Filter(filter OrderedDocument) DeleteOperation
	One() error
	Many() error
}

// Delete : returns a DeleteOperation instance
func (storage *MongoStorage) Delete(collectionName string) DeleteOperation {
	collection, ok := storage.collections[collectionName]
	op := MongoDeleteOperation{
		collection: collection,
		opts:       options.Delete(),
	}
	if !ok {
		op.err = errors.New("invalid collection: " + collectionName)
	}
	return op
}

// MongoDeleteOperation : implements DeleteOperation for MongoDB, it is used to delete one or many documents in the database
type MongoDeleteOperation struct {
	collection *mongo.Collection
	ctx        context.Context
	opts       *options.DeleteOptions
	filter     OrderedDocument
	err        error
}

// Context : sets the context for the operation
func (do MongoDeleteOperation) Context(ctx context.Context) DeleteOperation {
	do.ctx = ctx
	return do
}

// Filter : sets the filter for the operation
func (do MongoDeleteOperation) Filter(filter OrderedDocument) DeleteOperation {
	for _, elem := range filter {
		do.filter = append(do.filter, primitive.E{Key: elem.Key, Value: elem.Value})
	}

	return do
}

// One : deletes one document in the database
func (do MongoDeleteOperation) One() error {
	if do.err != nil {
		return do.err
	}

	result, err := do.collection.DeleteOne(do.ctx, do.filter, do.opts)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("nothing to delete")
	}

	return nil
}

// Many : deletes many documents in the database
func (do MongoDeleteOperation) Many() error {
	if do.err != nil {
		return do.err
	}

	result, err := do.collection.DeleteMany(do.ctx, do.filter, do.opts)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("nothing to delete")
	}

	return nil
}
