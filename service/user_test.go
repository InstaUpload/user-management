package service

import (
	"context"
	"errors"
	"testing"

	common "github.com/InstaUpload/common/types"
	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
)

func TestCreate(t *testing.T) {
	mockService := NewService(&store.MockStore)
	ctx := context.Background()
	user := types.CreateUserPayload{
		Name:     "Sahaj",
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	t.Run("Pass Create User", func(t *testing.T) {
		if err := mockService.User.Create(ctx, &user); err != nil {
			t.Errorf("Can not create user, err: %v", err)
		}
	})
	t.Run("Fail when existing user tries to create account", func(t *testing.T) {
		err := mockService.User.Create(ctx, &user)
		if !errors.Is(err, common.ErrDataFound) {
			t.Errorf("Expected ErrDataFound, but got %v", err)
		}
	})
}

func TestLogin(t *testing.T) {
	mockService := NewService(&store.MockStore)
	ctx := context.Background()
	user := types.LoginUserPayload{
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	t.Run("Pass Login User", func(t *testing.T) {
		token, err := mockService.User.Login(ctx, &user)
		if err != nil {
			t.Errorf("Can not login user, err: %v", err)
		}
		if token == "" {
			t.Errorf("Expected token but got empty string")
		}
	})
	user = types.LoginUserPayload{
		Email:    "gpt.sahaj@gmail.com",
		Password: "password123",
	}
	t.Run("Fail when non existing user tries to login", func(t *testing.T) {
		token, err := mockService.User.Login(ctx, &user)
		if !errors.Is(err, common.ErrDataNotFound) {
			t.Errorf("Expected ErrDataNotFound, but got %v", err)
		}
		if token != "" {
			t.Errorf("Expected empty string as token, but got %s", token)
		}
	})
}
