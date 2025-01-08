package utils

type AuthAction string

const (
	Register AuthAction = "create"
	Login    AuthAction = "login"
)

type Credential struct {
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Action   AuthAction `json:"action"`
}

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
