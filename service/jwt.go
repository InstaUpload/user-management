package service

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	jwtExpire time.Time
	jwtSecret []byte
}

func (j *JWTService) GenerateToken(userId int64) (string, error) {
	// convert userId to string.

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.FormatInt(userId, 10),
		"expiredAt": j.jwtExpire,
	})
	return token.SignedString(j.jwtSecret)
}

func (j *JWTService) ParseToken(token string) (int64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.jwtSecret, nil
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
