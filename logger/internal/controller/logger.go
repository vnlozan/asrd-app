package controller

import (
	"logger/internal/dto"
	"logger/internal/service"
	"logger/internal/utils"
	"net/http"
)

type LoggerController struct {
	loggerService service.ILoggerService
}

func NewLoggerController(loggerService service.ILoggerService) *LoggerController {
	return &LoggerController{
		loggerService: loggerService,
	}
}

func (c *LoggerController) AddOneLog(w http.ResponseWriter, r *http.Request) {
	var logEntry dto.LogEntry
	if err := utils.ReadJSON(w, r, &logEntry); err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	if err := c.loggerService.AddOneLog(logEntry); err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JsonResponse{Error: false, Message: "logged"}
	utils.WriteJSON(w, http.StatusAccepted, resp)
}
