package service

import (
	"context"
	"errors"
	"testing"

	common "github.com/InstaUpload/common/types"
	"github.com/InstaUpload/user-management/types"
)

func TestCreate(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	user := types.CreateUserPayload{
		Name:     "Sahaj",
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	t.Run("Pass Create User", func(t *testing.T) {
		if err := mockService.User.Create(testCtx, &user); err != nil {
			t.Errorf("Can not create user, err: %v", err)
		}
	})
	t.Run("Fail when existing user tries to create account", func(t *testing.T) {
		err := mockService.User.Create(testCtx, &user)
		if !errors.Is(err, common.ErrDataFound) {
			t.Errorf("Expected ErrDataFound, but got %v", err)
		}
	})
}

func TestLogin(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	user := types.LoginUserPayload{
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	t.Run("Pass Login User", func(t *testing.T) {
		token, err := mockService.User.Login(testCtx, &user)
		if err != nil {
			t.Errorf("Can not login user, err: %v", err)
		}
		if token == "" {
			t.Errorf("Expected token but got empty string")
		}
		testCtx = context.WithValue(testCtx, "TestLoginJWT", token)
	})
	user = types.LoginUserPayload{
		Email:    "gpt.sahaj@gmail.com",
		Password: "password123",
	}
	t.Run("Fail when non existing user tries to login", func(t *testing.T) {
		token, err := mockService.User.Login(testCtx, &user)
		if !errors.Is(err, common.ErrDataNotFound) {
			t.Errorf("Expected ErrDataNotFound, but got %v", err)
		}
		if token != "" {
			t.Errorf("Expected empty string as token, but got %s", token)
		}
	})
}

func TestAuth(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	token := testCtx.Value("TestLoginJWT").(string)
	t.Run("Pass user auth", func(t *testing.T) {
		user, err := mockService.User.Auth(testCtx, token)
		if err != nil {
			t.Errorf("Expected no error but found one: %v", err)
		}
		if user.Id == 0 {
			t.Errorf("Expected user id 1 but got 0")
		}
		testCtx = context.WithValue(testCtx, TestUser, user)
	})
	token = token + "invalid"
	t.Run("Fail user auth", func(t *testing.T) {
		user, err := mockService.User.Auth(testCtx, token)
		if !errors.Is(err, common.ErrDataNotFound) {
			t.Errorf("Expected Incorrect data error, but got %s", err.Error())
		}
		if user.Id != 0 {
			t.Errorf("Expected user id to be 0 but got %d", user.Id)
		}
	})
}

func TestUpdateRole(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	testAdminUser, ok := testCtx.Value(TestUser).(types.User)
	if !ok {
		t.Errorf("Need TestUser to perform test")
	}
	testAdminUser.Role.Name = "admin"
	testCtx = context.WithValue(testCtx, TestUser, testAdminUser)
	// NOTE: types.User is expected in ctx.
	t.Run("Pass User Update Role function", func(t *testing.T) {
		if err := mockService.User.UpdateRole(testCtx, 1, "admin"); err != nil {
			t.Errorf("Expected no error but got %v", err.Error())
		}
	})
	t.Run("Fail User Update Role function", func(t *testing.T) {
		if err := mockService.User.UpdateRole(testCtx, 1, "superAdmin"); err != nil {
			if !errors.Is(err, common.ErrDataNotFound) {
				t.Errorf("Expected data not found error but got %v", err.Error())
			}
		}
	})
	testAdminUser.Role.Name = "regular"
	testCtx = context.WithValue(testCtx, TestUser, testAdminUser)
}

func TestResetPassword(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	testAdminUser, ok := testCtx.Value(TestUser).(types.User)
	if !ok {
		t.Errorf("Need TestUser to perform test")
	}
	email := testAdminUser.Email
	t.Run("Reset password", func(t *testing.T) {
		token, err := mockService.User.ResetPassword(testCtx, email)
		if err != nil {
			t.Errorf("Expected no error but got %v", err.Error())
		}
		if token == "" {
			t.Errorf("Expected token to be set")
		}
		testCtx = context.WithValue(testCtx, TestPasswordToken, token)
	})
}
