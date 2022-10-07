package main

import (
	"dailypractice/tip"
	"dailypractice/user"
	"dailypractice/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getTips(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":    tip.All(),
	})
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":    user.All(),
	})
}

func login(c *gin.Context) {
	email := c.PostForm("email")
	plainPW := c.PostForm("password")

	user, err := user.FindByEmail(email)
  fmt.Printf("emai: %s, password: %s, user: %+v", email, plainPW, user)
	utils.CheckError(err)

	if user.AuthByPassword(plainPW) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "error",
		})
	}
}

func main() {
	router := gin.Default()
	router.GET("api/tips", getTips)
	router.GET("api/users", getUsers)
	router.POST("api/login", login)
  router.Run(":9000")
}
