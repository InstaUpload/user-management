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
	dbstore *store.Store
}

func (s *UserService) Create(ctx context.Context, userPayload *types.UserPayload) error {
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
	// convert user to userPayload.
	userPayload.Id = user.Id
	userPayload.IsVerified = user.IsVerified
	userPayload.CreatedAt = user.CreatedAt

	return nil
}
