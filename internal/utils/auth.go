package utils

import (
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"net/smtp"
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

func generateOTP() string {
	otp := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", otp)
}

func SendSMTPMail(to string, SMTPConfig ct.SMTPConfig, db *gorm.DB) error {
	verificationCode := generateOTP()

	var otp ct.OTP

	otp.Email = to
	otp.OTP = verificationCode

	err := db.Save(&otp).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	body := fmt.Sprintf("Your verification code is:\n %v", verificationCode)

	// Compose message
	message := []byte(fmt.Sprintf("From: %s\r\n", SMTPConfig.From) + fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, "Email Verification", body))

	// Authentication
	auth := smtp.PlainAuth("", SMTPConfig.Username, SMTPConfig.Password, SMTPConfig.Host)

	err = smtp.SendMail(SMTPConfig.Host+":"+SMTPConfig.Port, auth, SMTPConfig.From, []string{to}, message)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
