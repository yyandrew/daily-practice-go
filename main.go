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

type TipForm struct {
	Context  string `form:"context"`
	Category string `form:"category"`
	ImageUrl string `form:"imageUrl"`
}

func getTips(c *gin.Context) {
	var tips = make([]tip.Tip, 0)
	category := c.DefaultQuery("category", "vim")
	content := c.DefaultQuery("context", "")
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	result, err := user.FindById(user_id)
	if err == nil {
		currentUser := result.(user.User)
		userTips, err := currentUser.Tips(category, content)
		if err != nil {
			fmt.Println(err)
		} else {
			tips = userTips.(tip.Tipslice).Tips
		}
	} else {
		tips = tip.All(category).Tips
	}

	c.JSON(http.StatusOK, tips)
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
	form := TipForm{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, err := token.ExtractTokenID(c)
	utils.CheckError(err)
	fmt.Printf("user_id %s\n", user_id)
	result, err := user.FindById(user_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	currentUser := result.(user.User)

	result, err = tip.Create(form.Context, form.Category, form.ImageUrl, string(currentUser.Id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		newTip := result.(tip.Tip)
		c.JSON(http.StatusOK, newTip)
	}
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
	public.DELETE("/tips/:id", middlewares.JwtAuthMiddleware(), controllers.DeleteTip)
	router.Run(":9000")
}
