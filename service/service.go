package service

import (
	"context"

	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

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
