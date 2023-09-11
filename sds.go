package sds

import "context"

// Repo provides a way of storeing a generic Go item in a database.
type Repo[A any] interface {
	Find(ctx context.Context, id string) (A, error)
	Save(ctx context.Context, dh A) error
	Query(ctx context.Context, opts ...QueryOption) ([]A, Cursor, error)
	Delete(ctx context.Context, id string) error
}

// Cursor is a way to keep track of what to fetch next.
type Cursor = *string

// WithCursor adds a cursor to a query.
func WithCursor(cursor Cursor) QueryOption {
	return func(o *opt) {
		o.cursor = cursor
	}
}

// Desending returns the query results in desending order based on creation time.
func Descending() QueryOption {
	return func(o *opt) {
		o.descending = true
	}
}

// WithLimit limits the number of items returned in a query.
func WithLimit(n int) QueryOption {
	return func(o *opt) {
		o.limit = n
	}
}

// Entity is an interface that some implementations will need to support
// to index items in a database.
type Entity interface {
	GetID() string
}

// QueryOption provides a way to supply a cursor or a limit to a query.
type QueryOption = func(*opt)

type opt struct {
	cursor     *string
	limit      int
	descending bool
}

// Descending determines the order in which results should be returned.
func (o *opt) Descending() bool {
	return o.descending
}

// Limit determines how many items should be returned
func (o *opt) Limit() int {
	return o.limit
}

// Cursor determines where to begin returning results.
func (o *opt) Cursor() Cursor {
	return o.cursor
}

// Options is a set of query parameters.
type Options interface {
	Descending() bool
	Limit() int
	Cursor() Cursor
}

// MakeOpts returns a set of query parameters.
func MakeOpts(opts []QueryOption) Options {

	options := &opt{
		descending: false,
		limit:      25,
		cursor:     nil,
	}

	for _, f := range opts {
		f(options)
	}

	return options
}
