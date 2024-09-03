package routes

import (
	"net/http"

	"github.com/cobaltbase/cobaltbase/internal/cobaltbase/ct"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ApiRouter() *chi.Mux {
	type js = ct.Js
	api := chi.NewRouter()

	api.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, js{
			"message": "/api endpoint",
		})
	})

	api.Mount("/tables", TablesRouter())
	api.Mount("/items", ItemsRouter())

	return api
}
