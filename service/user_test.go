package service

import (
	"context"
	"errors"
	"testing"
	"time"

	common "github.com/InstaUpload/common/types"
	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
)

func TestCreate(t *testing.T) {
	mockService := NewService(&store.MockStore)
	ctx := context.Background()
	user := types.UserPayload{
		Name:     "Sahaj",
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	t.Run("Pass Create User", func(t *testing.T) {
		if err := mockService.User.Create(ctx, &user); err != nil {
			t.Errorf("Can not create user, err: %v", err)
		}
		if user.CreatedAt.Day() != time.Now().Day() {
			t.Errorf("Logic error in creating user got %v date want %v", user.CreatedAt.Day(), time.Now().Day())
		}
	})
	t.Run("Fail Create User", func(t *testing.T) {
		err := mockService.User.Create(ctx, &user)
		if !errors.Is(err, common.ErrDataFound) {
			t.Errorf("Expected ErrorDataFound, but got %v", err)
		}
	})
}

func TestLogin(t *testing.T) {
	mockService := NewService(&store.MockStore)
	ctx := context.Background()
	user := types.UserPayload{
		Name:     "Sahaj",
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	t.Run("Pass Login User", func(t *testing.T) {
	})
}
