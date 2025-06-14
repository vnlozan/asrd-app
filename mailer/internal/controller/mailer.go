package controller

import (
	"log"
	"mailer/internal/dto"
	"mailer/internal/service"
	"mailer/internal/utils"
	"net/http"
)

type MailController struct {
	MailerService service.IMailerService
}

func NewMailController(mailerService service.IMailerService) *MailController {
	return &MailController{
		MailerService: mailerService,
	}
}

func (m *MailController) SendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload dto.SendMailMessageRequest

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		utils.ErrorJSON(w, err)
		return
	}

	msg := dto.MailMessage{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = m.MailerService.SendMessage(msg)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	payload := utils.JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	utils.WriteJSON(w, http.StatusAccepted, payload)
}
