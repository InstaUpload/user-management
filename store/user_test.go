package store

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	log.Printf("In user test function")
	ctx := context.Background()
	user := &User{
		Name:  "Sahaj",
		Email: "gpt.sahaj28@gmail.com",
		Password: Password{
			text: "password123",
		},
	}
	t.Run("Pass Create User", func(t *testing.T) {
		if err := MockStore.User.Create(ctx, user); err != nil {
			log.Printf("%v", user)
			t.Errorf("Can not create user, err: %v", err)
		}
		if user.CreatedAt.Day() != time.Now().Day() {
			t.Errorf("user creation went wrong, got %d want %d", user.CreatedAt.Day(), time.Now().Day())
		}
	})
	// TODO: get mock store and test the user's create function.
}
