package handler

import (
	"github.com/dengzii/genx/testdata/param"
	"github.com/gin-gonic/gin"
)

//go:generate genx handler
func LoginHandler(ctx *gin.Context, request *param.LoginRequest) (*param.LoginResponse, error) {
	// ...
	return &param.LoginResponse{Token: "token"}, nil
}

//go:generate genx handler
func GetMessageHandler(ctx *gin.Context) ([]*param.Message, error) {
	// ...
	return []*param.Message{}, nil
}

//go:generate genx handler
func GetMessageHandler2(request *param.TestRequest) ([]*param.Message, error) {
	// ...
	return []*param.Message{}, nil
}
