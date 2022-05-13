package handler

import "github.com/dengzii/genx/testdata/gin"

type InfoController struct {
}

//go:generate genx api
func (*InfoController) GetInfo(context *gin.Context) error {
	return nil
}
