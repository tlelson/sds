package mem_test

import (
	"context"
	"testing"

	"github.com/schafer14/sds/mem"
	"github.com/schafer14/sds/test"
)

func TestMemRepo(t *testing.T) {

	t.Parallel()
	ctx := context.Background()
	store := mem.New[test.Entity]()

	// Can DoesItWork provide its own generic entity? Must be unique for a given
	// backend database?
	test.DoesItWork(t, ctx, store)

}
