package utils

import (
	"time"

	"github.com/cobaltbase/cobaltbase/internal/constants"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user ct.Js, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(constants.JWT_SECRET)
}
