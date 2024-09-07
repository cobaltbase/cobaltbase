package controllers

import (
	"net/http"
	"time"

	"github.com/cobaltbase/cobaltbase/internal/config"
	"github.com/cobaltbase/cobaltbase/internal/constants"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/cobaltbase/cobaltbase/internal/utils"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user ct.Auth
		var requestBody ct.AuthRequest

		if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
		}

		if err := config.Validate.Var(requestBody.Email, `email,required`); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
		}

		if err := config.Validate.Var(requestBody.Email, `required,min=8`); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
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

		user_map := ct.Js{"email": user.Email, "role": user.Role, "verified": user.Verified}

		// Generate JWT tokens
		accessToken, err := utils.GenerateJWT(user_map, 15*time.Minute)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateJWT(user_map, 15*24*time.Hour)
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
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		render.JSON(w, r, ct.Js{"message": "user authenticated successfully"})

	}
}

func ValidateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": "unauthorized"})
			return
		}

		tokenStr := cookie.Value
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": "unauthorized"})
			return
		}

		user := claims["user"]

		render.JSON(w, r, ct.Js{
			"message": "user is authorized",
			"user":    user,
		})
	}
}

func RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		refreshToken := cookie.Value

		var session ct.Session
		err = config.DB.First(&session, &ct.Session{RefreshToken: refreshToken}).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"error": "session not found"})
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(refreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			render.Status(r, 401)
			render.JSON(w, r, ct.Js{"error": "unauthorized"})
			return
		}

		user := claims["user"].(ct.Js)
		newAccessToken, err := utils.GenerateJWT(user, 15*time.Minute)
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, ct.Js{"error": "server error"})
			return
		}

		// Set new access token in HttpOnly cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			Expires:  time.Now().Add(15 * time.Minute),
			HttpOnly: true,
			Path:     "/",
		})

		w.WriteHeader(http.StatusOK)
	}
}
