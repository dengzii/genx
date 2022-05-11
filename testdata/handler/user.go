package handler

import (
	"genx/testdata/param"
	"github.com/gin-gonic/gin"
)

//go:generate genx api post /login
func LoginHandler(ctx *gin.Context, request *param.LoginRequest) (*param.LoginResponse, error) {
	// ...
	return &param.LoginResponse{Token: "token"}, nil
}

//go:generate genx api get /messages
func GetMessageHandler(ctx *gin.Context) ([]*param.Message, error) {
	// ...
	return []*param.Message{}, nil
}
