package main

import (
	"github.com/bitqiu/pix-gen/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	r.GET("/captcha", handler.HandleCaptcha)
	r.GET("/qrcode", handler.HandleQrcode)
	r.Run(":8080")
}
