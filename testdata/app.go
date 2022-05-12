package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.New()

	//engine.Handle("POST", "/login", handler.GenxLoginHandler)
	//engine.Handle("GET", "/msg", handler.GenxGetMessageHandler)

	err := engine.Run(":8080")

	if err != nil {
		panic(err)
	}
}
