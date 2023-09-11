package test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/schafer14/sds"
)

func DoesItWork[A sds.Entity](t *testing.T, ctx context.Context, s sds.Repo[A], save func(string) error) {

	list := []string{}
	reversed := []string{}

	for i := 0; i < 100; i++ {

		id := fmt.Sprintf("%.2d", i)
		list = append(list, id)
		reversed = append(reversed, id)
		err := save(id)
		if err != nil {
			t.Errorf("saving item %v : %v", id, err)
		}
	}

	sort.Slice(reversed, func(i, j int) bool {
		return reversed[j] < reversed[i]
	})

	t.Run("fetches items by id", func(_ *testing.T) {
		first, err := s.Find(ctx, "00")
		if err != nil {
			t.Errorf("finding item %v", err)
		}
		if list[0] != first.GetID() {
			t.Errorf("expected item %q, got %q", list[0], first.GetID())
		}

	})

	var curs sds.Cursor
	t.Run("fetches items with a limit", func(t *testing.T) {
		firstTwo, c, err := s.Query(ctx, sds.WithLimit(2))
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		curs = c
		if neq(list[:2], firstTwo) {
			t.Errorf("expected item %q, got %v", list[:2], firstTwo)
		}

	})

	var nextCurs sds.Cursor
	t.Run("fetches items with a cursor", func(t *testing.T) {
		nextTwo, c, err := s.Query(ctx, sds.WithLimit(2), sds.WithCursor(curs))
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		nextCurs = c
		if neq(list[2:4], nextTwo) {
			t.Errorf("expected item %q, got %v", list[2:4], nextTwo)
		}
	})

	var lastCur sds.Cursor
	t.Run("fetches items with a limit greater than the items in the db", func(t *testing.T) {
		lastNintySix, c, err := s.Query(ctx, sds.WithLimit(100), sds.WithCursor(nextCurs))
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		lastCur = c
		if len(lastNintySix) != 96 {
			t.Errorf("expected 96 items got %d", len(lastNintySix))
		}
		if lastCur != nil {
			t.Errorf("expecting cursor to be nil")
		}
		if neq(list[4:100], lastNintySix) {
			t.Errorf("expected item %q, got %v", list[4:100], lastNintySix)
		}
	})

	t.Run("fetch an empty list when a cursor has no more items left", func(t *testing.T) {
		_, zeroCur, err := s.Query(ctx, sds.WithLimit(96), sds.WithCursor(nextCurs))
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		empty, curless, err := s.Query(ctx, sds.WithCursor(zeroCur))
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		if curless != nil {
			t.Errorf("expected a nil cursor")
		}
		if len(empty) != 0 {
			t.Errorf("expected [ ] got %v", empty)
		}
	})

	var rcurs sds.Cursor
	t.Run("fetches items with a limit (descending)", func(t *testing.T) {
		firstTwo, c, err := s.Query(ctx, sds.WithLimit(2), sds.Descending())
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		rcurs = c
		if neq(reversed[:2], firstTwo) {
			t.Errorf("expected item %q, got %v", reversed[:2], firstTwo)
		}
	})

	var rnextCurs sds.Cursor
	t.Run("fetches items with a cursor (descending)", func(t *testing.T) {
		nextTwo, c, err := s.Query(ctx, sds.WithLimit(2), sds.WithCursor(rcurs), sds.Descending())
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		rnextCurs = c
		if neq(reversed[2:4], nextTwo) {
			t.Errorf("expected item %q, got %v", reversed[2:4], nextTwo)
		}
	})

	var rlastCur sds.Cursor
	t.Run("fetches items with a limit greater than the items in the db (descending)", func(t *testing.T) {
		lastNintySix, c, err := s.Query(ctx, sds.WithLimit(100), sds.WithCursor(rnextCurs), sds.Descending())
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		rlastCur = c
		if len(lastNintySix) != 96 {
			t.Errorf("expected 96 items got %d", len(lastNintySix))
		}
		if rlastCur != nil {
			t.Errorf("expecting cursor to be nil")
		}
		if neq(reversed[4:100], lastNintySix) {
			t.Errorf("expected item %q, got %v", reversed[4:100], lastNintySix)
		}
	})

	t.Run("fetch an empty list when a cursor has no more items left (descending)", func(t *testing.T) {
		_, zeroCur, err := s.Query(ctx, sds.WithLimit(96), sds.WithCursor(rnextCurs), sds.Descending())
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		empty, curless, err := s.Query(ctx, sds.WithCursor(zeroCur), sds.Descending())
		if err != nil {
			t.Errorf("querying items %v", err)
		}
		if curless != nil {
			t.Errorf("expected a nil cursor")
		}
		if len(empty) != 0 {
			t.Errorf("expected [ ] got %v", empty)
		}
	})

	t.Run("deletes items from the database", func(t *testing.T) {

		err := s.Delete(ctx, list[4])
		if err != nil {
			t.Errorf("deleting item %v", err)
		}
		_, err = s.Find(ctx, list[4])
		if err == nil {
			t.Errorf("item should not be returned after being deleted")
		}
		queryRes, _, err := s.Query(ctx, sds.WithLimit(100))
		if err != nil {
			t.Errorf("querying items : %v", err)
		}
		if len(queryRes) != 99 {
			fmt.Println(queryRes)
			t.Errorf("deleted items should not be returned in query")
		}

	})

}

func neq[A sds.Entity](as []string, bs []A) bool {
	return !eq(as, bs)
}

func eq[A sds.Entity](as []string, bs []A) bool {
	if len(as) != len(bs) {
		return false
	}

	for i := range as {
		if as[i] != bs[i].GetID() {
			return false
		}
	}

	return true
}
