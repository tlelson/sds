package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"

	"github.com/schafer14/sds"
)

type service[A sds.Entity] struct {
	db *mongo.Collection
}

// New creates a new data holder service.
func New[A sds.Entity](db *mongo.Collection) (sds.Repo[A], error) {

	return &service[A]{
		db: db,
	}, nil
}

func (s *service[A]) Find(ctx context.Context, id string) (A, error) {
	var item A
	err := s.db.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&item)
	return item, err
}

func (s *service[A]) Save(ctx context.Context, item A) error {
	if _, err := s.db.InsertOne(ctx, item); err != nil {
		return errors.Wrap(err, "saving item")
	}

	return nil
}

func (s *service[A]) Query(ctx context.Context, opts ...sds.QueryOption) ([]A, sds.Cursor, error) {

	opt := sds.MakeOpts(opts)
	filter := bson.D{{}}

	order := 1
	cursorDir := "$gt"
	if opt.Descending() {
		order = -1
		cursorDir = "$lt"
	}

	if opt.Cursor() != nil {
		filter = bson.D{{"_id", bson.D{{cursorDir, opt.Cursor()}}}}
	}

	sortOpts := options.Find().SetSort(bson.D{{"_id", order}}).SetLimit(int64(opt.Limit()))

	cursor, err := s.db.Find(ctx, filter, sortOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "finding items")
	}
	defer cursor.Close(ctx)

	result := []A{}
	for cursor.Next(ctx) {
		var item A
		err := cursor.Decode(&item)
		if err != nil {
			return nil, nil, errors.Wrap(err, "decoding item")
		}

		result = append(result, item)
	}

	var curs sds.Cursor
	if len(result) == opt.Limit() {
		c := result[len(result)-1].GetID()
		curs = &c
	}

	return result, curs, nil
}

func (s *service[A]) Delete(ctx context.Context, id string) error {

	_, err := s.db.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}
