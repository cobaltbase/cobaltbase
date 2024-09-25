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

type CobaltBase struct {
	Router *chi.Mux
}

var server http.Server

func New() *CobaltBase {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://accounts.google.com"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	router.Mount("/api", routes.ApiRouter())
	config.Configure()
	return &CobaltBase{
		Router: router,
	}
}

func (cb *CobaltBase) Run(port string) {
	log.Printf("Server is running on http://%s", port)
	server.Handler = cb.Router
	server.Addr = port
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
