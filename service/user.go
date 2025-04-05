package service

import (
	"context"
	"log"
	"strings"

	common "github.com/InstaUpload/common/types"
	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
)

type UserService struct {
	dbstore    *store.Store
	jwtService interface {
		GenerateToken(int64) (string, error)
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
	token, err := s.jwtService.GenerateToken(user.Id)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return "", common.ErrDataNotFound
	}
	// return token.
	return token, nil
}
