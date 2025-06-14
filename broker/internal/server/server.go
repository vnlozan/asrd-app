package server

import (
	"broker/internal/controller"
	"fmt"
	"log"
	"net/http"

	"broker/internal/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Config           *config.Config
	brokerController *controller.BrokerController
}

func NewServer(config *config.Config, brokerController *controller.BrokerController) *Server {
	return &Server{Config: config, brokerController: brokerController}
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

	mux.Post("/", s.brokerController.HandleRootRequest)
	mux.Post("/handle", s.brokerController.HandleRequest)

	return mux
}
