package middlewares

import (
	"context"
	"github.com/cobaltbase/cobaltbase/internal/config"
	"github.com/cobaltbase/cobaltbase/internal/constants"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/cobaltbase/cobaltbase/internal/utils"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string
		cookie, err := r.Cookie("access_token")
		if err != nil {
			refCookie, err := r.Cookie("refresh_token")
			if err != nil {
				render.Status(r, 401)
				render.JSON(w, r, ct.Js{"error": "unauthorized"})
				return
			}

			refreshToken := refCookie.Value

			newAccessToken, err := utils.RefreshToken(refreshToken, config.DB)
			if err != nil {
				cookieNames := []string{"access_token", "refresh_token"}

				// Get the current time
				now := time.Now()

				// Create and set expired cookies
				for _, name := range cookieNames {
					cookie := &http.Cookie{
						Name:     name,
						Value:    "",
						Path:     "/",
						HttpOnly: true,                    // Use this if your site is served over HTTPS
						Expires:  now.Add(-1 * time.Hour), // Set expiration to the past
						MaxAge:   -1,
					}
					http.SetCookie(w, cookie)
				}
				render.Status(r, 401)
				render.JSON(w, r, ct.Js{"error": "unauthorized"})
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    newAccessToken,
				Expires:  time.Now().Add(15 * time.Minute),
				HttpOnly: true,
				Path:     "/",
			})
			tokenStr = newAccessToken
		} else {
			tokenStr = cookie.Value
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return constants.JWT_SECRET, nil
		})

		if err != nil || !token.Valid {
			cookie, err := r.Cookie("refresh_token")
			if err != nil {
				cookieNames := []string{"access_token", "refresh_token"}

				// Get the current time
				now := time.Now()

				// Create and set expired cookies
				for _, name := range cookieNames {
					cookie := &http.Cookie{
						Name:     name,
						Value:    "",
						Path:     "/",
						HttpOnly: true,                    // Use this if your site is served over HTTPS
						Expires:  now.Add(-1 * time.Hour), // Set expiration to the past
						MaxAge:   -1,
					}
					http.SetCookie(w, cookie)
				}
				render.Status(r, 401)
				render.JSON(w, r, ct.Js{"error": "unauthorized"})
				return
			}

			refreshToken := cookie.Value

			newAccessToken, err := utils.RefreshToken(refreshToken, config.DB)
			if err != nil {
				render.Status(r, 401)
				render.JSON(w, r, ct.Js{"error": "unauthorized"})
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    newAccessToken,
				Expires:  time.Now().Add(15 * time.Minute),
				HttpOnly: true,
				Path:     "/",
			})
		}

		user := claims["user"]

		user, ok := user.(ct.Js)
		if !ok {
			render.Status(r, 401)
			render.JSON(w, r, ct.Js{"error": "unauthorized"})
			return
		}

		ctx := context.WithValue(r.Context(), ct.AuthMiddlewareKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
