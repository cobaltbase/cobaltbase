package routes

import (
	"github.com/go-chi/chi/v5"
)

func ApiRouter() *chi.Mux {
	api := chi.NewRouter()

	api.Mount("/tables", TablesRouter())
	api.Mount("/items", ItemsRouter())
	api.Mount("/auth", AuthRouter())
	api.Mount("/config", ConfigRouter())

	return api
}
