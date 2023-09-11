## Getting started

To get started you need to 

1. Create a type for the data you would like to store.
2. Implement the `sds.Entity` interface for your type.
3. Initiate a service.
4. Save/Find/Query/Delete your data.

Here is an example:

```go
package main

type user struct {
	ID     string `bson:"_id"`
	Name   string `bson:"name"`
	Email  string `bson:"email"` 
}

func (u user) GetID() string { return u.ID }

func main() {

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
}
```

This is a simple example to get you started. For more in-depth thoughts on 
how this library can facilitate domain driven development read [this article](./wiki/how-tos/domain-service.md).
