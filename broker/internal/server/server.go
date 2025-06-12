package server

import (
	"broker/internal/dto"
	"broker/internal/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"broker/internal/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Config config.Config
}

func NewServer(config config.Config) *Server {
	return &Server{Config: config}
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

func (s *Server) Broker(w http.ResponseWriter, r *http.Request) {
	payload := utils.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = utils.WriteJSON(w, http.StatusOK, payload)
}

func (s *Server) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload dto.RequestPayload

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		s.authenticate(w, requestPayload.Auth)
	case "log":
		s.logItem(w, requestPayload.Log)
	case "mail":
		s.sendMail(w, requestPayload.Mail)
	default:
		utils.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (s *Server) logItem(w http.ResponseWriter, entry dto.LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := fmt.Sprintf("%s/log", s.Config.LoggerURL)

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		utils.ErrorJSON(w, err)
		return
	}

	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	utils.WriteJSON(w, http.StatusAccepted, payload)

}

func (s *Server) authenticate(w http.ResponseWriter, a dto.AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/authenticate", s.Config.AuthURL), bytes.NewBuffer(jsonData))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		utils.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		utils.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService utils.JsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		utils.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	utils.WriteJSON(w, http.StatusAccepted, payload)
}

func (s *Server) sendMail(w http.ResponseWriter, msg dto.MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := fmt.Sprintf("%s/send", s.Config.MailerURL)

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		utils.ErrorJSON(w, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	utils.WriteJSON(w, http.StatusAccepted, payload)

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

	mux.Post("/", s.Broker)

	mux.Post("/handle", s.HandleSubmission)

	return mux
}
