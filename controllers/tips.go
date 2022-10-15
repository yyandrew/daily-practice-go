package controllers

import (
	"dailypractice/tip"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteTip(c *gin.Context) {
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
