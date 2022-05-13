package handler

import (
	"github.com/dengzii/genx/testdata/gin"
	"github.com/dengzii/genx/testdata/param"
)

//go:generate genx handler
func LoginHandler(ctx *gin.Context, request *param.LoginRequest) (*param.LoginResponse, error) {
	// ...
	return &param.LoginResponse{Token: "token"}, nil
}

//go:generate genx handler
func GetMessageHandler(ctx *gin.Context) {
	// ...

}

//go:generate genx handler
func GetMessageHandler2(request *param.TestRequest) ([]*param.Message, error) {
	// ...
	return []*param.Message{}, nil
}
