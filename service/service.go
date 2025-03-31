package service

import (
	"context"

	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
)

type Service struct {
	User interface {
		Create(context.Context, *types.UserPayload) error
	}
}

func NewService(dbstore *store.Store) Service {
	return Service{
		User: &UserService{
			dbstore,
		},
	}
}
