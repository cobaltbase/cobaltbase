package routes

import (
	"github.com/cobaltbase/cobaltbase/internal/controllers"
	"github.com/cobaltbase/cobaltbase/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func ItemsRouter() *chi.Mux {
	tr := chi.NewRouter()

	tr.Route("/{table}", func(r chi.Router) {
		r.Use(middlewares.CheckTableExists) // Attach the middleware here

		r.Get("/", controllers.GetAllItems())
		r.Get("/single", controllers.GetItem())
		r.With(middlewares.PreProcessingMiddleware).Post("/", controllers.CreateItem())
		r.With(middlewares.PreProcessingMiddleware).Put("/", controllers.UpdateItem())
		r.Delete("/{id}", controllers.DeleteItem())
	})

	return tr
}
