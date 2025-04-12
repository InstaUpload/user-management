package service

import (
	"context"
	"time"

	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
	"github.com/InstaUpload/user-management/utils"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Service struct {
	User interface {
		Create(context.Context, *types.CreateUserPayload) error
		Login(context.Context, *types.LoginUserPayload) (string, error)
		Auth(context.Context, string) (types.User, error)
		UpdateRole(context.Context, int64, string) error
	}
}

func NewService(dbstore *store.Store) Service {
	expTime := time.Now().Add(time.Hour * time.Duration(utils.GetEnvInt("JWTEXPHR", 24))).Unix()
	jwtService := &JWTService{
		jwtExpire: time.Unix(expTime, 0),
		jwtSecret: []byte(utils.GetEnvString("JWTSECRET", "secret")),
	}
	return Service{
		User: &UserService{
			dbstore,
			jwtService,
		},
	}
}
