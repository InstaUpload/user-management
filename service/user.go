package service

import (
	"context"
	"log"
	"slices"
	"strings"

	common "github.com/InstaUpload/common/types"
	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
)

const CurrentUser string = "CurrentUser"

// TODO: Create a const array name it superadmin and add emails in it.
var superadmin = []string{"gpt.sahaj@gmail.com"}

type UserService struct {
	dbstore    *store.Store
	jwtService interface {
		GenerateAuthToken(int64) (string, error)
		ParseAuthToken(string) (int64, error)
		GeneratePasswordToken(int64) (string, error)
		ParsePasswordToken(string) (int64, error)
	}
}

func (s *UserService) Create(ctx context.Context, userPayload *types.CreateUserPayload) error {
	// Add validation for userPayload struct.
	if err := validate.Struct(userPayload); err != nil {
		// return a custome error from errors package. for invalide data.
		log.Printf("invalid user data: %v", err)
		return common.ErrIncorrectDataReceived
	}
	// below section needs to be automated.
	var user types.User
	user.Name = userPayload.Name
	user.Email = userPayload.Email
	user.Password.Text = userPayload.Password

	user.Password.HashPassword()
	if err := s.dbstore.User.Create(ctx, &user); err != nil {
		// Check error, if already exist return common.ErrDataFound error.
		// write a if block that check if "duplicate key" error, return common.ErrDataFound error.
		if strings.Contains(err.Error(), "duplicate key") {
			return common.ErrDataFound
		}
		log.Printf("err: %s", err.Error())
		return err
	}
	// TODO: Send varification email to user.

	return nil
}

func (s *UserService) Login(ctx context.Context, userPayload *types.LoginUserPayload) (string, error) {
	// convert userPayload to user.
	user := types.User{
		Email: userPayload.Email,
	}
	// call user.GetUserByEmail() method.
	// check if error is returned from Login().
	// if error is not nil, return common.ErrDataNotFound error.
	if err := s.dbstore.User.GetUserByEmail(ctx, &user); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", common.ErrDataNotFound
		}

		log.Printf("err: %s", err.Error())
		return "", err
	}
	// match user.Password.hash with userPayload.Password.
	// If match fails return common.ErrIncorrectDataReceived error.
	user.Password.Text = userPayload.Password
	if err := user.Password.ComparePassword(); err != nil {
		return "", common.ErrIncorrectDataReceived
	}
	// generate token and return it.
	token, err := s.jwtService.GenerateAuthToken(user.Id)
	if err != nil {
		log.Printf("err: %s", err.Error())
		// TODO: To be updated to internal server error.
		return "", common.ErrDataNotFound
	}
	// return token.
	return token, nil
}

func (s *UserService) Auth(ctx context.Context, token string) (types.User, error) {
	// get userId from token.
	var user types.User
	userId, err := s.jwtService.ParseAuthToken(token)
	if err != nil {
		// check the error message and return error accordingly.
		if strings.Contains(err.Error(), "token is expired") {
			return user, common.ErrIncorrectDataReceived
		} else if strings.Contains(err.Error(), "token is invalid") {
			return user, common.ErrDataNotFound
		}
	}
	// get user by id.
	user.Id = userId
	// call user.GetUserById() method.
	// if error is not nil, return common.ErrDataNotFound error.
	if err := s.dbstore.User.GetUserById(ctx, &user); err != nil {
		// check if err contains.
		if strings.Contains(err.Error(), "no rows") {
			return user, common.ErrDataNotFound
		}
		log.Printf("err: %s", err.Error())
		return user, err
	}
	// return user.
	return user, nil
}

func (s *UserService) UpdateRole(ctx context.Context, userId int64, roleName string) error {
	// Check if current user has admin role. or current user is one of super admin users.
	currentUser := ctx.Value(CurrentUser).(types.User)
	if currentUser.Role.Name != "admin" {
		return common.ErrIncorrectDataReceived
	}
	// Get user by passed id.
	var userToBeUpdated types.User
	userToBeUpdated.Id = userId
	if err := s.dbstore.User.GetUserById(ctx, &userToBeUpdated); err != nil {
		return common.ErrIncorrectDataReceived
	}
	// Check if user to be updated in super admin user, if so return an Incorrect data error.
	if slices.Contains(superadmin, userToBeUpdated.Email) {
		return common.ErrIncorrectDataReceived
	}
	// Update user role to passed roleName.

	if err := s.dbstore.User.UpdateUserRole(ctx, &userToBeUpdated, roleName); err != nil {
		// TODO: check the error type if it containes role deosnt exist.
		// Then return ErrIncorrectDataReceived error.
		if strings.Contains(err.Error(), "no rows") {
			return common.ErrDataNotFound
		}
		log.Printf("err: %s", err.Error())
		return err

	}
	// Return nil if update is successful.
	return nil
}
