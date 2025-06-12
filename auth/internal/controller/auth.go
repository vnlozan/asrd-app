package controller

import (
	"auth/internal/dto"
	"auth/internal/service"
	"auth/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthController struct {
	authService service.IAuthService
}

func NewAuthController(authService service.IAuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (a *AuthController) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload dto.AuthRequest

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := a.authService.Authenticate(r.Context(), requestPayload)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	payload := utils.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	err = logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, payload)
}

func logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger:8080/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
