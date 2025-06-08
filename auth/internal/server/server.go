package server

import (
	"auth/internal/controller"
	"fmt"
	"log"
	"net/http"

	"auth/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Config         config.Config
	AuthController controller.AuthController
}

func NewServer(config config.Config, authController controller.AuthController) *Server {
	return &Server{Config: config, AuthController: authController}
}

func (s *Server) Start() {
	log.Println("Starting on port", s.Config.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Config.Port),
		Handler: s.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (s *Server) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/authenticate", s.AuthController.Authenticate)
	return mux
}
