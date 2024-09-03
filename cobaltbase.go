package cobaltbase

import (
	"log"
	"net/http"

	"github.com/cobaltbase/cobaltbase/internal/config"
	"github.com/cobaltbase/cobaltbase/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type cobaltBase struct {
	Router *chi.Mux
}

func New() *cobaltBase {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Mount("/api", routes.ApiRouter())
	config.Configure()
	return &cobaltBase{
		Router: router,
	}
}

func (cb *cobaltBase) Run(port string) {
	log.Printf("Server is running on %s", port)
	if err := http.ListenAndServe(port, cb.Router); err != nil {
		log.Fatal(err)
	}
}
