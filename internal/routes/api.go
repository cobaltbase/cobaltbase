package routes

import (
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/go-chi/chi/v5"
)

func ApiRouter() *chi.Mux {
	type js = ct.Js
	api := chi.NewRouter()

	api.Mount("/tables", TablesRouter())
	api.Mount("/items", ItemsRouter())
	api.Mount("/auth", AuthRouter())

	return api
}
