package main

import (
	"dailypractice/controllers"
	"dailypractice/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Static("public", "./public")
	public := router.Group("/api")
	public.GET("/tips", controllers.GetTips)
	public.POST("/signup", controllers.SignUp)
	public.POST("/login", controllers.Login)
	router.POST("api/upload", controllers.Upload)
	public.POST("/tips", middlewares.JwtAuthMiddleware(), controllers.CreateTip)
	public.DELETE("/tips/:id", middlewares.JwtAuthMiddleware(), controllers.DeleteTip)
	router.Run(":9000")
}
