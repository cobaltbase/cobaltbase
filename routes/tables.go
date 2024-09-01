package routes

import (
	"github.com/cobaltbase/cobaltbase/controllers"
	"github.com/go-chi/chi/v5"
)

func TablesRouter() *chi.Mux {
	tr := chi.NewRouter()

	tr.Get("/", controllers.GetAllTables())
	tr.Get("/{table}/schema", controllers.GetSchema())

	tr.Post("/", controllers.CreateTable())
	tr.Post("/{table}/field", controllers.UpdateSingleField())
	tr.Delete("/{table}/field", controllers.DeleteSingleField())

	return tr
}
