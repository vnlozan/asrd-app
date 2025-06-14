package dto

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type SendMailMessageRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type MailMessage struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}
