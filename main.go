package main

import (
	"dailypractice/tip"
	"dailypractice/user"
	"dailypractice/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getTips(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": tip.All(),
	})
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": user.All(),
	})
}

func upload(c *gin.Context) {
	fileParam := "tip-image"
	file, _ := c.FormFile(fileParam)
	dst := "./public/img/" + fileParam + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(file.Filename)
	err := c.SaveUploadedFile(file, dst)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"filepath": filepath.Base(dst),
		})
	} else {
		fmt.Println(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "upload error",
		})
	}
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
			"user":    map[string]string{"email": user.Email},
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "error",
		})
	}
}

func main() {
	router := gin.Default()
	router.Static("public", "./public")
	router.GET("api/tips", getTips)
	router.GET("api/users", getUsers)
	router.POST("api/login", login)
	router.POST("api/upload", upload)
	router.Run(":9000")
}
