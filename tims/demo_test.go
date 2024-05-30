package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/schafer14/sds/mem"
	"github.com/segmentio/ksuid"
)

type user struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

func (u user) GetID() string { return u.ID }

func TestDemo(t *testing.T) {
	is := is.New(t)

	// Setup the storage repository
	ctx := context.Background()
	userRepo := mem.New[user]()

	// Save an item
	id := ksuid.New().String()
	err := userRepo.Save(ctx, user{
		ID:    id,
		Name:  "Banner",
		Email: "banner@example.com",
	})
	is.NoErr(err)

	// Retrieve the item
	user, err := userRepo.Find(ctx, id)
	is.NoErr(err)

	// Check everything worked
	is.Equal(user.Email, "banner@example.com")
	is.Equal(user.Name, "Banner")
	is.Equal(user.ID, id)
	fmt.Printf("%#v\n", user)
}
