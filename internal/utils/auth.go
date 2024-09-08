package utils

import (
	"gorm.io/gorm"
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
func RefreshToken(refreshToken string, db *gorm.DB) (string, error) {
	var session ct.Session
	err := db.First(&session, &ct.Session{RefreshToken: refreshToken}).Error
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return constants.JWT_SECRET, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	user := claims["user"].(ct.Js)
	newAccessToken, err := GenerateJWT(user, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}
