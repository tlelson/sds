package bbolt

import (
	"context"
	"math"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/schafer14/sds"
)

type bboltRepo[A sds.Entity] struct {
	db         *bbolt.DB
	bucketName []byte
}

func New[A sds.Entity](db *bbolt.DB, bucketName string) (sds.Repo[A], error) {

	bucket := []byte(bucketName)
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})

	return &bboltRepo[A]{
		db:         db,
		bucketName: bucket,
	}, err
}

func (b *bboltRepo[A]) Find(ctx context.Context, id string) (A, error) {

	var a A

	err := b.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket([]byte(b.bucketName)).Get([]byte(id))

		return bson.Unmarshal(value, &a)
	})
	if err != nil {
		return a, err
	}

	return a, nil
}

func (b *bboltRepo[A]) Save(ctx context.Context, item A) error {

	data, err := bson.Marshal(item)
	if err != nil {
		return err
	}

	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucketName)

		return bucket.Put([]byte(item.GetID()), data)
	})
}

func (b *bboltRepo[A]) Query(ctx context.Context, opts ...sds.QueryOption) ([]A, sds.Cursor, error) {

	options := sds.MakeOpts(opts)

	var result []A
	err := b.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(b.bucketName).Cursor()
		var setup func() ([]byte, []byte)

		if options.Cursor() == nil && !options.Descending() {
			setup = func() ([]byte, []byte) {
				return c.Seek([]byte(""))
			}
		} else if options.Cursor() == nil {
			setup = func() ([]byte, []byte) {
				return c.Last()
			}
		} else if !options.Descending() {
			setup = func() ([]byte, []byte) {
				curs := *options.Cursor()
				c.Seek([]byte(curs))
				return c.Next()
			}
		} else {
			setup = func() ([]byte, []byte) {
				curs := *options.Cursor()
				c.Seek([]byte(curs))
				return c.Prev()
			}
		}

		it := c.Next
		if options.Descending() {
			it = c.Prev
		}

		for k, v := setup(); k != nil && options.Limit() > len(result); k, v = it() {
			var i A
			err := bson.Unmarshal(v, &i)
			if err != nil {
				return err
			}

			result = append(result, i)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	length := int(math.Min(float64(options.Limit()), float64(len(result))))
	var cursor sds.Cursor

	if length == options.Limit() {
		c := result[length-1].GetID()
		cursor = &c
	}

	return result, cursor, nil

}

func (b *bboltRepo[A]) Delete(ctx context.Context, id string) error {

	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucketName)

		return bucket.Delete([]byte(id))
	})
}
