package cobaltbase

import (
	"log"
	"net/http"

	"github.com/cobaltbase/cobaltbase/internal/cobaltbase/config"
	"github.com/cobaltbase/cobaltbase/internal/cobaltbase/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type cobaltBase struct {
	Router *chi.Mux
}

// newCobaltBase creates a new instance of the internal cobaltBase.
func New() *cobaltBase {
	router := chi.NewRouter()
	return &cobaltBase{
		Router: router,
	}
}

func (cb *cobaltBase) Run(port string) {
	cb.Router.Use(middleware.Logger)
	cb.Router.Mount("/api", routes.ApiRouter())
	config.Configure()
	log.Printf("Server is running on %s", port)
	if err := http.ListenAndServe(port, cb.Router); err != nil {
		log.Fatal(err)
	}
}
