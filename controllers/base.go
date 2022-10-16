package controllers

import (
	. "dailypractice/utils/constants"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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
