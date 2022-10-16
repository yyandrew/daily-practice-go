package controllers

import (
	"dailypractice/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var IMG_PATH string

func init() {
	err := godotenv.Load()
	utils.CheckError(err)
	IMG_PATH = os.Getenv("IMG_PATH")
}

func Upload(c *gin.Context) {
	fileParam := "tip-image"
	file, _ := c.FormFile(fileParam)
	dst := IMG_PATH + fileParam + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(file.Filename)
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
