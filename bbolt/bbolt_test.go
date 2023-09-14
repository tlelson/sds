package bbolt_test

import (
	"context"
	"path"
	"testing"

	"github.com/matryer/is"
	"github.com/schafer14/sds"
	bboltStorage "github.com/schafer14/sds/bbolt"
	"github.com/schafer14/sds/test"
	"github.com/segmentio/ksuid"
	"go.etcd.io/bbolt"
)

func TestBboltDB(t *testing.T) {

	t.Parallel()
	ctx := context.Background()
	is := is.New(t)
	dir := t.TempDir()
	p := path.Join(dir, "test.db")
	db, err := bbolt.Open(p, 0600, nil)
	is.NoErr(err)
	defer db.Close()

	store, err := bboltStorage.New[test.Entity](db, ksuid.New().String())
	is.NoErr(err)

	test.DoesItWork(t, ctx, store)

}

func TestBboltDBDataStructure(t *testing.T) {

	t.Parallel()
	ctx := context.Background()
	is := is.New(t)
	dir := t.TempDir()
	p := path.Join(dir, "test.db")
	db, err := bbolt.Open(p, 0600, nil)
	is.NoErr(err)
	defer db.Close()

	store, err := bboltStorage.New[test.Entity](db, ksuid.New().String())
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
