package service

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// modify these attributes divide them into 4, 2 for login 2 for update password.
type JWTService struct {
	authExpire     time.Time
	authSecret     []byte
	passwordExpire time.Time
	passwordSecret []byte
}

func (j *JWTService) GenerateAuthToken(userId int64) (string, error) {
	// convert userId to string.

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.FormatInt(userId, 10),
		"expiredAt": j.authExpire,
	})
	return token.SignedString(j.authSecret)
}

func (j *JWTService) ParseAuthToken(token string) (int64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.authSecret, nil
	})
	if err != nil {
		return 0, err
	}
	strUserId := claims["userId"].(string)
	userId, err := strconv.ParseInt(strUserId, 10, 64)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (j *JWTService) GeneratePasswordToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.FormatInt(userId, 10),
		"expiredAt": j.passwordExpire,
	})
	return token.SignedString(j.passwordSecret)
}

func (j *JWTService) ParsePasswordToken(token string) (int64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.passwordSecret, nil
	})
	if err != nil {
		return 0, err
	}
	strUserId := claims["userId"].(string)
	userId, err := strconv.ParseInt(strUserId, 10, 64)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
