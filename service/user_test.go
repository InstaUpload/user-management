package service

import (
	"context"
	"errors"
	"testing"

	common "github.com/InstaUpload/common/types"
	"github.com/InstaUpload/user-management/types"
)

var testUsers = []struct {
	Name           string
	Email          string
	Password       string
	AuthToken      string
	PasswordToken  string
	VerifyToken    string
	EditorReqToken string
	ctx            context.Context
}{
	{
		Name:     "Sahaj 1",
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
		ctx:      context.Background(),
	},
	{
		Name:     "Sahaj 2",
		Email:    "gpt.sahaj2@gmail.com",
		Password: "password456",
		ctx:      context.Background(),
	},
}

// NOTE: user id for test users are 1 and 3.
// NOTE: Where userId 1 is verified and authenticated user and userId 3 is not verified and unauthenticated user.

func TestCreate(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}

	for u := range testUsers {
		user := types.CreateUserPayload{
			Name:     testUsers[u].Name,
			Email:    testUsers[u].Email,
			Password: testUsers[u].Password,
		}
		t.Run("Pass Create User", func(t *testing.T) {
			if err := mockService.User.Create(testUsers[u].ctx, &user); err != nil {
				t.Errorf("Can not create user, err: %v", err)
			}
		})
		t.Run("Fail when existing user tries to create account", func(t *testing.T) {
			err := mockService.User.Create(testUsers[u].ctx, &user)
			if !errors.Is(err, common.ErrDataFound) {
				t.Errorf("Expected ErrDataFound, but got %v", err)
			}
		})
	}

}

func TestLogin(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	for u := range testUsers {
		user := types.LoginUserPayload{
			Email:    testUsers[u].Email,
			Password: testUsers[u].Password,
		}

		t.Run("Pass Login User", func(t *testing.T) {
			token, err := mockService.User.Login(testUsers[u].ctx, &user)
			if err != nil {
				t.Errorf("Can not login user, err: %v", err)
			}
			if token == "" {
				t.Errorf("Expected token but got empty string")
			}
			testUsers[u].AuthToken = token
		})
		user.Email = "not@exist.com"
		t.Run("Fail when non existing user tries to login", func(t *testing.T) {
			token, err := mockService.User.Login(testUsers[u].ctx, &user)
			if !errors.Is(err, common.ErrDataNotFound) {
				t.Errorf("Expected ErrDataNotFound, but got %v", err)
			}
			if token != "" {
				t.Errorf("Expected empty string as token, but got %s", token)
			}
		})
	}
}

func TestAuth(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	for u := range testUsers {
		token := testUsers[u].AuthToken
		t.Run("Pass user auth", func(t *testing.T) {
			user, err := mockService.User.Auth(testUsers[u].ctx, token)
			if err != nil {
				t.Errorf("Expected no error but found one: %v", err)
			}
			if user.Id == 0 {
				t.Errorf("Expected user id %d but got %d", user.Id, 0)
			}
			tempCtx := testUsers[u].ctx
			testUsers[u].ctx = context.WithValue(tempCtx, common.CurrentUserKey, user)
		})
		token = token + "invalid"
		t.Run("Fail user auth", func(t *testing.T) {
			user, err := mockService.User.Auth(testUsers[u].ctx, token)
			if !errors.Is(err, common.ErrDataNotFound) {
				t.Errorf("Expected Incorrect data error, but got %s", err.Error())
			}
			if user.Id != 0 {
				t.Errorf("Expected user id to be 0 but got %d", user.Id)
			}
		})
	}
}

func TestUpdateRole(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	// NOTE: types.User is expected in ctx.
	t.Run("Pass User Update Role function", func(t *testing.T) {
		if err := mockService.User.UpdateRole(testUsers[0].ctx, 3, "admin"); err != nil {
			t.Errorf("Expected no error but got %v", err.Error())
		}
	})
	t.Run("Fail User Update Role function", func(t *testing.T) {
		if err := mockService.User.UpdateRole(testUsers[0].ctx, 3, "superAdmin"); err != nil {
			if !errors.Is(err, common.ErrDataNotFound) {
				t.Errorf("Expected data not found error but got %v", err.Error())
			}
		}
	})
}

func TestResetPassword(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	for u := range testUsers {
		t.Run("Reset password", func(t *testing.T) {
			token, err := mockService.User.ResetPassword(testUsers[u].ctx, testUsers[u].Email)
			if err != nil {
				t.Errorf("Expected no error but got %v", err.Error())
			}
			if token == "" {
				t.Errorf("Expected token to be set")
			}
			testUsers[u].PasswordToken = token
		})

		t.Run("Fail Reset password", func(t *testing.T) {
			email := "notfound@gmail.com"
			token, err := mockService.User.ResetPassword(testUsers[u].ctx, email)
			if !errors.Is(err, common.ErrIncorrectDataReceived) {
				t.Errorf("Expected Incorrect data error, but got %s", err.Error())
			}
			if token != "" {
				t.Errorf("Did not expect token but got %s", token)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	password := "updated password"
	for u := range testUsers {
		t.Run("Pass Update user password using token", func(t *testing.T) {
			if err := mockService.User.UpdatePassword(testUsers[u].ctx, testUsers[u].PasswordToken, password); err != nil {
				t.Errorf("Did not expect error but got %v", err)
			}
			user := types.LoginUserPayload{
				Email:    "gpt.sahaj28@gmail.com",
				Password: password,
			}
			token, err := mockService.User.Login(testUsers[u].ctx, &user)
			if err != nil {
				t.Errorf("Did not expect error but got %v", err)
			}
			if token == "" {
				t.Error("Expected token but got none.")
			}
		})

		t.Run("Fail Update user passr", func(t *testing.T) {
			token := testUsers[u].PasswordToken + "invalid"
			if err := mockService.User.UpdatePassword(testUsers[u].ctx, token, "<PASSWORD>"); err != nil {
				if !errors.Is(err, common.ErrDataNotFound) {
					t.Errorf("Expected data not for error but got %v", err)
				}
			}
		})
	}
}

func TestSendVerification(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	// NOTE: types.User is expected in ctx.
	t.Run("Pass Send varification to user function", func(t *testing.T) {
		token, err := mockService.User.SendVerification(testUsers[0].ctx)
		if err != nil {
			t.Errorf("did not expect error but got %v", err)
		} else if token == "" {
			t.Errorf("Expected token to be not empty")
		}
		testUsers[0].VerifyToken = token
	})
}

// NOTE: From this point forward testUsers[0] is verified user and testUsers[1] is not verified user.
func TestVerify(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	token := testUsers[0].VerifyToken
	t.Run("Pass Verify user function", func(t *testing.T) {
		if err := mockService.User.Verify(testCtx, token); err != nil {
			t.Errorf("did not expect error but got %v", err)
		}
	})
	t.Run("Fail Verify user function with Verified user", func(t *testing.T) {
		if err := mockService.User.Verify(testCtx, token); err != nil {
			if !errors.Is(err, common.ErrIncorrectDataReceived) {
				t.Errorf("Expected Incorrect data recevied error but got %v", err)
			}
		}
	})
	token = token + "Invalid"
	t.Run("Fail Verify user function with Invalid token", func(t *testing.T) {
		if err := mockService.User.Verify(testCtx, token); err != nil {
			if !errors.Is(err, common.ErrDataNotFound) {
				t.Errorf("Expected data not found error but got %v", err)
			}
		}
	})
}

func TestSendEditorRequest(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	var userId int64 = 3
	t.Run("Pass Send editor request function", func(t *testing.T) {
		res, err := mockService.User.SendEditorRequest(testUsers[0].ctx, userId)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if res.Token == "" {
			t.Errorf("Expected token but got empty string")
		}
		testUsers[0].EditorReqToken = res.Token
	})
	userId = 2
	t.Run("Fail Send editor request function with non existing user", func(t *testing.T) {
		_, err := mockService.User.SendEditorRequest(testUsers[0].ctx, userId)
		if !errors.Is(err, common.ErrDataNotFound) {
			t.Errorf("Expected data not found error but got %v", err)
		}
	})
}

func TestAddEditor(t *testing.T) {
	mockService, ok := testCtx.Value(MockService).(Service)
	if !ok {
		t.Errorf("Need MockService to perform test")
	}
	// NOTE: In this function I need editor in the ctx
	// NOTE: CurrentUser is expected in token.
	t.Run("Pass Add editor function", func(t *testing.T) {
		if err := mockService.User.AddEditor(testUsers[1].ctx, testUsers[0].EditorReqToken); err != nil {
			t.Errorf("Expected no error but found %v", err)
		}
	})
	t.Run("Fail add editor funtion on adding already added editor", func(t *testing.T) {
		if err := mockService.User.AddEditor(testUsers[1].ctx, testUsers[0].EditorReqToken); err != nil {
			if !errors.Is(err, common.ErrIncorrectDataReceived) {
				t.Errorf("Expected Incorrect data recevied but got %v", err)
			}
		}
	})
}
