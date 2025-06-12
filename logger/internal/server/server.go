package server

import (
	"fmt"
	"log"
	"logger/internal/config"
	"logger/internal/controller"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Config           config.Config
	LoggerController controller.LoggerController
}

func NewServer(config *config.Config, loggerController *controller.LoggerController) *Server {
	return &Server{
		Config:           *config,
		LoggerController: *loggerController,
	}
}

func (s *Server) Start() {
	log.Println("Starting service on port", s.Config.WebPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Config.WebPort),
		Handler: s.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
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
	mux.Post("/log", s.LoggerController.WriteLog)

	return mux
}
