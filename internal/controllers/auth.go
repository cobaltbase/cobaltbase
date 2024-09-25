package controllers

import (
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"net/http"
	"time"

	"github.com/cobaltbase/cobaltbase/internal/config"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/cobaltbase/cobaltbase/internal/utils"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user ct.Auth
		var requestBody ct.AuthRequest

		if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
			return
		}

		if err := config.Validate.Var(requestBody.Email, `email,required`); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
			return
		}

		if err := config.Validate.Var(requestBody.Email, `required,min=8`); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		user.Email = requestBody.Email
		user.Password = string(hashedPassword)
		user.Verified = false

		if err := config.DB.Create(&user).Error; err != nil {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		user.Password = ""
		render.JSON(w, r, ct.Js{"message": "User Created Succesfully", "user": user})
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ct.AuthRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		var user ct.Auth
		if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		userMap := ct.Js{"email": user.Email, "role": user.Role, "verified": user.Verified, "id": user.ID}

		// Generate JWT tokens
		accessToken, err := utils.GenerateJWT(userMap, 15*time.Minute)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateJWT(userMap, 15*24*time.Hour)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var session ct.Session

		session.AuthID = user.ID
		session.Provider = "local"
		session.UserAgent = r.Header.Get("User-Agent")
		session.RefreshToken = refreshToken

		err = config.DB.Create(&session).Error
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": "could not create session"})
			return
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

		render.JSON(w, r, ct.Js{"message": "user authenticated successfully"})

	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refCookie, err := r.Cookie("refresh_token")
		if err != nil {
			render.Status(r, 401)
			render.JSON(w, r, ct.Js{"error": "unauthorized"})
			return
		}

		refreshToken := refCookie.Value

		if err := config.DB.Unscoped().Delete(&ct.Session{}, &ct.Session{
			RefreshToken: refreshToken,
		}).Error; err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "could not delete session"})
			return
		}

		cookie := &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,                           // Use this if your site is served over HTTPS
			Expires:  time.Now().Add(-1 * time.Hour), // Set expiration to the past
			MaxAge:   -1,
		}
		http.SetCookie(w, cookie)
		cookie = &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,                           // Use this if your site is served over HTTPS
			Expires:  time.Now().Add(-1 * time.Hour), // Set expiration to the past
			MaxAge:   -1,
		}
		http.SetCookie(w, cookie)

		render.JSON(w, r, ct.Js{"message": "user logged out"})
	}
}

func ValidateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(ct.AuthMiddlewareKey).(ct.Js)

		if !ok {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": "internal server error"})
			return
		}

		render.JSON(w, r, ct.Js{
			"message": "user is authorized",
			"user":    user,
		})
	}
}

func GetSessions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessions []ct.Session

		user, ok := (r.Context().Value(ct.AuthMiddlewareKey)).(ct.Js)

		if !ok {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": "internal server error"})
			return
		}

		err := config.DB.Select("id", "created_at", "updated_at", "user_agent", "auth_id", "provider").Find(&sessions, &ct.Session{
			AuthID: user["id"].(string),
		}).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		render.JSON(w, r, ct.Js{"sessions": sessions})

	}
}

func RevokeSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := (r.Context().Value(ct.AuthMiddlewareKey)).(ct.Js)

		if !ok {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": "internal server error"})
			return
		}

		var input struct {
			ID string `json:"id"`
		}
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		err := config.DB.Where("id = ?", input.ID).Delete(&ct.Session{}, &ct.Session{
			AuthID: user["id"].(string),
		}).Error

		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		render.JSON(w, r, ct.Js{"message": "session revoked successfully"})
	}
}

func SendMailCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Email string `json:"email" validate:"required,email"`
		}

		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		err = config.Validate.Struct(user)
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		go func() {
			err := utils.SendSMTPMail(user.Email, config.SMTPConfig, config.DB)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Mail sent")
			}
		}()
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		render.JSON(w, r, ct.Js{"message": "mail will be sent"})

	}
}

func VerifyEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ct.AuthMiddlewareKey).(ct.Js)
		if !ok {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": "internal server error"})
			return
		}
		var otp ct.OTP
		var userDb ct.Auth
		var otpReq struct {
			OTP string `json:"otp"`
		}

		if err := render.DecodeJSON(r.Body, &otpReq); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		if err := config.DB.First(&otp, &ct.OTP{
			Email: user["email"].(string),
		}).Error; err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": "otp not found"})
			return
		}

		if err := config.DB.First(&userDb, &ct.Auth{
			Email: user["email"].(string),
		}).Error; err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": "user not found"})
			return
		}

		if otp.OTP != otpReq.OTP {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": "wrong otp"})
			return
		}

		if !time.Now().Before(otp.UpdatedAt.Add(15 * time.Minute)) {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": "otp expired"})
			return
		}

		userDb.Verified = true
		config.DB.Save(&userDb)
		render.JSON(w, r, ct.Js{"message": "user verified successfully"})

	}
}

func ResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			OTP         string `json:"otp"`
			NewPassword string `json:"new_password"`
			OldPassword string `json:"old_password"`
			Email       string `json:"email" validate:"required,email"`
		}

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		if err := config.Validate.Struct(req); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		var user ct.Auth
		if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if req.OldPassword != "" {

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
		}

		if req.OTP == "" {

			var otp ct.OTP
			if err := config.DB.First(&otp, &ct.OTP{
				Email: user.Email,
			}).Error; err != nil {
				render.Status(r, 400)
				render.JSON(w, r, ct.Js{"error": "otp not found"})
				return
			}

			if otp.OTP != req.OTP {
				render.Status(r, 400)
				render.JSON(w, r, ct.Js{"error": "wrong otp"})
				return
			}

			if !time.Now().Before(otp.UpdatedAt.Add(15 * time.Minute)) {
				render.Status(r, 400)
				render.JSON(w, r, ct.Js{"error": "otp expired"})
				return
			}
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		user.Password = string(hashedPassword)

		if err := config.DB.Save(&user).Error; err != nil {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		//revoke all sessions
		if err := config.DB.Unscoped().Delete(&ct.Session{}, &ct.Session{
			AuthID: user.ID,
		}).Error; err != nil {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}

		userMap := ct.Js{"email": user.Email, "role": user.Role, "verified": user.Verified, "id": user.ID}

		// Generate JWT tokens
		accessToken, err := utils.GenerateJWT(userMap, 15*time.Minute)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateJWT(userMap, 15*24*time.Hour)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var session ct.Session

		session.AuthID = user.ID
		session.Provider = "local"
		session.UserAgent = r.Header.Get("User-Agent")
		session.RefreshToken = refreshToken

		err = config.DB.Create(&session).Error
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": "could not create session"})
			return
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

		render.JSON(w, r, ct.Js{"message": "password changed, all sessions revoked and current session logged in"})
	}
}

func ProviderAuthCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user goth.User
		var err error
		user, err = gothic.CompleteUserAuth(w, r)
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		err = utils.CompleteProviderAuth(user, w, r, config.DB)
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
	}
}

func ProviderAuthLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if user, err := gothic.CompleteUserAuth(w, r); err == nil {
			err := utils.CompleteProviderAuth(user, w, r, config.DB)
			if err != nil {
				render.Status(r, 400)
				render.JSON(w, r, ct.Js{"error": err.Error()})
				return
			}
		}
		url, err := gothic.GetAuthURL(w, r)
		if err != nil {
			render.Status(r, 401)
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		render.JSON(w, r, url)
	}
}

func CookieTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "test",
			Value:    "ajkshgfidsagkjlfgskg",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
		render.JSON(w, r, ct.Js{"message": "test logged in"})
	}
}
