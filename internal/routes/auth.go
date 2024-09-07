package routes

import (
	"github.com/cobaltbase/cobaltbase/internal/controllers"
	"github.com/go-chi/chi/v5"
)

func AuthRouter() *chi.Mux {
	ar := chi.NewRouter()

	ar.Post("/register", controllers.RegisterUser())
	ar.Post("/login", controllers.Login())

	ar.Get("/validate", controllers.ValidateToken())

	ar.Get("/refresh_token", controllers.RefreshToken())

	return ar
}
