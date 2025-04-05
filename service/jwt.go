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
