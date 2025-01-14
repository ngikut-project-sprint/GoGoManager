package models

type AuthAction string

const (
	Register AuthAction = "create"
	Login    AuthAction = "login"
)

type AuthRequest struct {
	Email    string     `json:"email" validate:"required"`
	Password string     `json:"password" validate:"required,min=8,max=32"`
	Action   AuthAction `json:"action"`
}

type AuthResponse struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
