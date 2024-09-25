package utils

import (
	"errors"
	"fmt"
	"github.com/markbates/goth"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
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
	var user ct.Auth
	if err := db.Where("email = ?", claims["user"].(ct.Js)["email"]).First(&user).Error; err != nil {
		return "", err
	}
	userMap := ct.Js{"email": user.Email, "role": user.Role, "verified": user.Verified, "id": user.ID}

	//user := claims["user"].(ct.Js)
	newAccessToken, err := GenerateJWT(userMap, 15*time.Minute)
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
	message := []byte(fmt.Sprintf("From: %s<%s>\r\n", SMTPConfig.FromName, SMTPConfig.From) + fmt.Sprintf("To: %s\r\n"+
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

func CompleteProviderAuth(user goth.User, w http.ResponseWriter, r *http.Request, db *gorm.DB) error {
	var dbUser ct.Auth
	if err := db.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dbUser.Email = user.Email
			dbUser.Verified = false
			if err := db.Save(&dbUser).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	userMap := ct.Js{"email": dbUser.Email, "role": dbUser.Role, "verified": dbUser.Verified, "id": dbUser.ID}

	// Generate JWT tokens
	accessToken, err := GenerateJWT(userMap, 15*time.Minute)
	if err != nil {
		return err
	}
	refreshToken, err := GenerateJWT(userMap, 15*24*time.Hour)
	if err != nil {
		return err
	}

	var session ct.Session

	session.AuthID = dbUser.ID
	session.Provider = "local"
	session.UserAgent = r.Header.Get("User-Agent")
	session.RefreshToken = refreshToken

	err = db.Create(&session).Error
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(15 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
            <html>
                <body>
                    <script>
                        window.close();
                    </script>
                    <p>Closing tab...</p>
                </body>
            </html>
        `))
	return nil
}
