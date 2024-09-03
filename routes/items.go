package routes

import (
	"github.com/cobaltbase/cobaltbase/controllers"
	"github.com/cobaltbase/cobaltbase/middlewares"
	"github.com/go-chi/chi/v5"
)

func ItemsRouter() *chi.Mux {
	tr := chi.NewRouter()

	tr.Get("/{table}", controllers.GetAllItems())

	tr.Get("/{table}/single", controllers.GetItem())
	tr.With(middlewares.PreProcessingMiddleware).Post("/{table}", controllers.CreateItem())
	tr.With(middlewares.PreProcessingMiddleware).Put("/{table}", controllers.UpdateItem())
	tr.Delete("/{table}/{id}", controllers.DeleteItem())

	return tr
}
