package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/schafer14/sds"
	mongoStorage "github.com/schafer14/sds/mongo"
	"github.com/schafer14/sds/test"
)

func TestMongoDB(t *testing.T) {

	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	is := is.New(t)
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cancel()
		if err := client.Disconnect(context.TODO()); err != nil {
			t.Fatal(err)
		}
	}()
	coll := client.Database("test").Collection("test_" + ksuid.New().String())

	store, err := mongoStorage.New[test.Entity](coll)
	is.NoErr(err)

	test.DoesItWork(t, ctx, store)

}

func TestMongoDBDataStructure(t *testing.T) {

	t.Parallel()
	ctx := context.Background()
	is := is.New(t)
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			t.Fatal(err)
		}
	}()
	coll := client.Database("test").Collection("test_" + ksuid.New().String())

	store, err := mongoStorage.New[test.Entity](coll)
	is.NoErr(err)

	err = store.Save(ctx, test.Entity{
		ID:    "abc",
		Field: "123",
	})
	is.NoErr(err)

	res, err := store.Find(ctx, "abc")
	is.NoErr(err)

	is.Equal(res.Field, "123")
	is.Equal(res.ID, "abc")

	lots, curs, err := store.Query(ctx, sds.WithLimit(42))
	is.NoErr(err)
	is.Equal(curs, nil)
	is.Equal(len(lots), 1)
	is.Equal(lots[0].Field, "123")
	is.Equal(lots[0].ID, "abc")
}
