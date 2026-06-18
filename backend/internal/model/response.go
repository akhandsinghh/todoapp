package model

type AuthResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}
type MessageResponse struct {
	Message string `json:"message"`
}
