package controller

import (
	"log"
	"mailer/internal/dto"
	"mailer/internal/repo/client"
	"mailer/internal/utils"
	"net/http"
)

type MailController struct {
	Mailer client.IMailerClient
}

func (m *MailController) SendMail(w http.ResponseWriter, r *http.Request) {

	var requestPayload dto.MailMessage

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		utils.ErrorJSON(w, err)
		return
	}

	msg := dto.Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = m.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		utils.ErrorJSON(w, err)
		return
	}

	payload := utils.JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	utils.WriteJSON(w, http.StatusAccepted, payload)
}
