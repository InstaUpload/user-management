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
		ResetPassword(context.Context, string) (string, error)
	}
}

func NewService(dbstore *store.Store) Service {
	// TODO: below lines needs to be update to add time and secret for password reset.
	expTime := time.Now().Add(time.Hour * time.Duration(utils.GetEnvInt("JWTEXPHR", 24))).Unix()
	passwordExpTime := time.Now().Add(time.Second * time.Duration(utils.GetEnvInt("JWTPASSWORDEXPTIME", 240))).Unix()
	jwtService := &JWTService{
		authExpire:     time.Unix(expTime, 0),
		authSecret:     []byte(utils.GetEnvString("JWTSECRET", "secret")),
		passwordExpire: time.Unix(passwordExpTime, 0),
		passwordSecret: []byte(utils.GetEnvString("JWTPASSWORDEXPTIME", "secret")),
	}
	return Service{
		User: &UserService{
			dbstore,
			jwtService,
		},
	}
}
