package main

import (
	"dailypractice/config"
	"dailypractice/controllers"
	"dailypractice/middlewares"
	. "dailypractice/utils/constants"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cnf, err := config.NewParsedConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	router := gin.Default()
	fmt.Printf("IMG_PATH: %s\n", IMG_PATH)

	router.Static("public", "./public")
	public := router.Group("/api")
	public.GET("/tips", controllers.GetTips)
	public.POST("/signup", controllers.SignUp)
	public.POST("/login", controllers.Login)
	public.POST("/logout", controllers.Logout)
	router.POST("api/upload", controllers.Upload)
	public.POST("/tips", middlewares.JwtAuthMiddleware(), controllers.CreateTip)
	public.DELETE("/tips/:id", middlewares.JwtAuthMiddleware(), controllers.DeleteTip)
	router.Run(fmt.Sprintf(":%d", cnf.ServerPort))
}
