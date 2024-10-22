package routes

import (
	"github.com/cobaltbase/cobaltbase/internal/controllers"
	"github.com/cobaltbase/cobaltbase/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func AuthRouter() *chi.Mux {
	ar := chi.NewRouter()

	ar.Post("/register", controllers.RegisterUser())
	ar.Post("/login", controllers.Login())
	ar.Delete("/logout", controllers.Logout())

	ar.With(middlewares.AuthenticateUser).Get("/validate", controllers.ValidateToken())
	ar.With(middlewares.AuthenticateUser).Get("/sessions", controllers.GetSessions())
	ar.With(middlewares.AuthenticateUser).Delete("/session", controllers.RevokeSession())
	ar.Post("/send-verification-mail", controllers.SendMailCode())
	ar.With(middlewares.AuthenticateUser).Post("/verify-email", controllers.VerifyEmail())
	ar.Post("/reset-password", controllers.ResetPassword())
	ar.Get("/oauth/callback", controllers.ProviderAuthCallback())
	ar.Get("/oauth/login", controllers.ProviderAuthLogin())
	ar.Get("/test-cookie", controllers.CookieTest())

	return ar
}
