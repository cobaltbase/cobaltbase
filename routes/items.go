package routes

import (
	"github.com/cobaltbase/cobaltbase/controllers"
	"github.com/cobaltbase/cobaltbase/middlewares"
	"github.com/go-chi/chi/v5"
)

func ItemsRouter() *chi.Mux {
	tr := chi.NewRouter()

	tr.Get("/{table}", controllers.GetAllItems())
	tr.With(middlewares.PreProcessingMiddleware).Post("/{table}", controllers.CreateItem())
	return tr
}
