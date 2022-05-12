package param

type LoginResponse struct {
	Token string
}

type LoginRequest struct {
	Name     string
	Password string
}

type Message struct {
	Content string
}

type TestRequest struct {
	Name string
}
