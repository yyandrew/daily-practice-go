package main

import (
	"dailypractice/controllers"
	"dailypractice/middlewares"
	"dailypractice/tip"
	"dailypractice/user"
	"dailypractice/utils"
	"dailypractice/utils/token"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getTips(c *gin.Context) {
	category := c.DefaultQuery("category", "vim")
	c.JSON(http.StatusOK, tip.All(category).Tips)
}

func deleteTip(c *gin.Context) {
	id := c.Param("id")
	res, ok := tip.Delete(id)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"tip":     res,
		})
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "unable to deleted tip",
		})
	}
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

func createTip(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)
	utils.CheckError(err)
	fmt.Printf("user_id %s\n", user_id)
	fmt.Printf("create tip\n")
}

func main() {
	router := gin.Default()

	router.Static("public", "./public")
	public := router.Group("/api")
	public.GET("/tips", getTips)
	public.GET("/users", getUsers)
	public.POST("/signup", controllers.SignUp)
	public.POST("/login", controllers.Login)
	router.POST("api/upload", upload)
	public.POST("/tips", middlewares.JwtAuthMiddleware(), createTip)
	router.Run(":9000")
}
