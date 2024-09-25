package routes

import (
	"github.com/cobaltbase/cobaltbase/internal/controllers"
	"github.com/cobaltbase/cobaltbase/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func ConfigRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/oauth", controllers.CreateOauthConfig())
	router.Get("/oauth/{provider}", controllers.RetrieveOauthConfig())
	router.Get("/oauth", controllers.ListOauthConfig())
	router.Put("/oauth", controllers.UpdateOauthConfig())
	router.Delete("/oauth", controllers.RemoveOauthConfig())

	router.With(middlewares.AuthenticateUser).Post("/smtp", controllers.UpdateSMTPConfig())

	return router
}
