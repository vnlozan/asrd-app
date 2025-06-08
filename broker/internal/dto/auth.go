package dto

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
