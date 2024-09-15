package cobaltbase

import (
	"github.com/go-chi/cors"
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
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	router.Mount("/api", routes.ApiRouter())
	config.Configure()
	return &cobaltBase{
		Router: router,
	}
}

func (cb *cobaltBase) Run(port string) {
	log.Printf("Server is running on http://%s", port)
	if err := http.ListenAndServe(port, cb.Router); err != nil {
		log.Fatal(err)
	}
}
