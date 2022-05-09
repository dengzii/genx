package param

type LoginResponse struct {
	Token string
}

type LoginRequest struct {
	Name     string
	Password string
}
