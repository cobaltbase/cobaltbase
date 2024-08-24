package routes

import (
	"github.com/cobaltbase/cobaltbase/controllers"
	"github.com/go-chi/chi/v5"
)

func TablesRouter() *chi.Mux {
	tr := chi.NewRouter()

	tr.Get("/", controllers.GetAllTables())
	tr.Post("/", controllers.CreateTable())

	return tr
}
