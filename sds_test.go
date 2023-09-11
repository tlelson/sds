package sds_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/segmentio/ksuid"

	"github.com/schafer14/sds/mem"
)

type user struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

func (u user) GetID() string { return u.ID }

func TestCreatingAReposito(t *testing.T) {

	ctx := context.Background()
	is := is.New(t)

	// Setup the repository.
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

}
