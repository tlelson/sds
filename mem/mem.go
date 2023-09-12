package mem

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/schafer14/sds"
)

type memRepo[A sds.Entity] struct {
	lock  *sync.RWMutex
	items map[string]A
}

func New[A sds.Entity]() sds.Repo[A] {
	return &memRepo[A]{
		lock:  &sync.RWMutex{},
		items: map[string]A{},
	}

}

func (f *memRepo[A]) Find(ctx context.Context, id string) (A, error) {

	f.lock.RLock()
	defer f.lock.RUnlock()

	item, ok := f.items[id]
	if !ok {
		return item, fmt.Errorf("item not found")
	}

	return item, nil
}

func (f *memRepo[A]) Query(ctx context.Context, opts ...sds.QueryOption) ([]A, sds.Cursor, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	options := sds.MakeOpts(opts)

	values := []A{}
	for _, item := range f.items {
		if options.Cursor() == nil {
			values = append(values, item)
		} else if item.GetID() < *options.Cursor() && options.Descending() {
			// Did you deliberately use the `GetID` function here rather than the key of the
			// f.items map ?
			values = append(values, item)
		} else if item.GetID() > *options.Cursor() && !options.Descending() {
			values = append(values, item)
		}
	}

	sort.Slice(values, func(i, j int) bool {
		if options.Descending() {
			return values[i].GetID() > values[j].GetID()
		}
		return values[i].GetID() < values[j].GetID()
	})

	length := int(math.Min(float64(options.Limit()), float64(len(values))))
	var cursor sds.Cursor

	if length == options.Limit() {
		c := values[length-1].GetID()
		cursor = &c
	}

	return values[:length], cursor, nil
}

func (f *memRepo[A]) Save(ctx context.Context, item A) error {
	f.lock.RLock()
	defer f.lock.RUnlock()

	f.items[item.GetID()] = item
	return nil
}

func (f *memRepo[A]) Delete(ctx context.Context, id string) error {
	f.lock.RLock()
	defer f.lock.RUnlock()

	delete(f.items, id)

	return nil
}
