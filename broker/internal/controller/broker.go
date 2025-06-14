package controller

import (
	"broker/internal/dto"
	"broker/internal/service"
	"broker/internal/utils"
	"errors"
	"log"
	"net/http"
)

type BrokerController struct {
	brokerService service.IBrokerService
}

func NewBrokerController(brokerService service.IBrokerService) *BrokerController {
	return &BrokerController{brokerService: brokerService}
}

func (c *BrokerController) HandleRootRequest(w http.ResponseWriter, r *http.Request) {
	payload := utils.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = utils.WriteJSON(w, http.StatusOK, payload)
}

func (c *BrokerController) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var requestPayload dto.RequestPayload

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	log.Println(requestPayload)

	switch requestPayload.Action {
	case "auth":
		{
			data, err := c.brokerService.Authenticate(requestPayload.Auth)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}

			payload := utils.JsonResponse{
				Error:   false,
				Message: "Authenticated!",
				Data:    data,
			}

			utils.WriteJSON(w, http.StatusAccepted, payload)
		}
	case "log":
		{
			err := c.brokerService.LogItem(requestPayload.Log)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}

			payload := utils.JsonResponse{
				Error:   false,
				Message: "logged via RabbitMQ",
			}

			utils.WriteJSON(w, http.StatusAccepted, payload)
		}
	case "mail":
		{
			err := c.brokerService.SendMail(requestPayload.Mail)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}

			payload := utils.JsonResponse{
				Error:   false,
				Message: "Message sent to " + requestPayload.Mail.To,
			}
			utils.WriteJSON(w, http.StatusAccepted, payload)
		}
	default:
		{
			utils.ErrorJSON(w, errors.New("unknown action"))
		}
	}
}
