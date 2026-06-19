package model

type AuthResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}
type MessageResponse struct {
	Message string `json:"message"`
}
type TaskListResponse struct {
	Items []TaskDTO `json:"items"`
	Total int64     `json:"total"`
}
