package routes

import (
	"github.com/cobaltbase/cobaltbase/internal/controllers"
	"github.com/go-chi/chi/v5"
)

func TablesRouter() *chi.Mux {
	tr := chi.NewRouter()

	tr.Get("/", controllers.GetAllTables())
	tr.Post("/", controllers.CreateTable())

	tr.Get("/{table}/schema", controllers.GetSchema())
	tr.Post("/{table}/field", controllers.UpdateSingleField())
	tr.Delete("/{table}/field", controllers.DeleteSingleField())

	tr.Delete("/{table}", controllers.DeleteSchemaWithoutData())
	tr.Delete("/{table}/WithDataDanger", controllers.DeleteSchemaWithData())

	return tr
}
