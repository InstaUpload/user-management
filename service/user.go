package service

import (
	"context"

	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
)

type UserService struct {
	dbstore *store.Store
}

func (s *UserService) Create(ctx context.Context, userPayload *types.UserPayload) error {
	// below section needs to be automated.
	var user types.User
	user.Id = userPayload.Id
	user.Name = userPayload.Name
	user.Email = userPayload.Email
	user.Password.Text = userPayload.Password
	user.IsVerified = userPayload.IsVerified
	user.CreatedAt = userPayload.CreatedAt

	user.Password.HashPassword()
	s.dbstore.User.Create(ctx, &user)
	// TODO: Send varification email to user.
	return nil
}
