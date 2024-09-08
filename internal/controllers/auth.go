package controllers

import (
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
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(15 * 24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		render.JSON(w, r, ct.Js{"message": "user authenticated successfully"})

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
