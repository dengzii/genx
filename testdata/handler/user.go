package handler

import (
	"genx/testdata/gin"
	"genx/testdata/param"
)

//go:generate genx api
//genx:api
func LoginHandler(ctx *gin.Context, request *param.LoginRequest, response *param.LoginResponse) error {

	return nil
}
