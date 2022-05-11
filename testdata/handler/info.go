package handler

import "github.com/gin-gonic/gin"

type InfoController struct {
}

//go:generate genx api
func (*InfoController) GetInfo(context *gin.Context) error {
	return nil
}
