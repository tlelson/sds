package mem_test

import (
	"context"
	"testing"

	"github.com/schafer14/sds/mem"
	"github.com/schafer14/sds/test"
)

type someEnt struct {
	id string
}

func (s *someEnt) GetID() string {
	return s.id
}

func (s *someEnt) String() string {
	return s.id
}

func TestMemRepo(t *testing.T) {

	t.Parallel()
	ctx := context.Background()
	store := mem.New[*someEnt]()

	// Can DoesItWork provide its own generic entity? Must be unique for a given
	// backend database?
	test.DoesItWork(t, ctx, store, func(s string) error {
		return store.Save(ctx, &someEnt{s})

	})

}
