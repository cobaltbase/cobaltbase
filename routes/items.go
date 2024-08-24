package routes

import (
	"github.com/cobaltbase/cobaltbase/controllers"
	"github.com/go-chi/chi/v5"
)

func ItemsRouter() *chi.Mux {
	tr := chi.NewRouter()

	tr.Get("/{table}", controllers.GetAllItems())
	return tr
}
