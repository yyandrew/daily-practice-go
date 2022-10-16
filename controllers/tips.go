package controllers

import (
	"dailypractice/tip"
	"dailypractice/user"
	"dailypractice/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TipForm struct {
	Context  string `form:"context"`
	Category string `form:"category"`
	ImageUrl string `form:"imageUrl"`
}

func GetTips(c *gin.Context) {
	var tips = make([]tip.Tip, 0)
	category := c.DefaultQuery("category", "vim")
	content := c.DefaultQuery("context", "")
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := user.FindById(user_id)
	if err == nil {
		currentUser := result.(user.User)
		userTips, err := currentUser.Tips(category, content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			tips = userTips.(tip.Tipslice).Tips
		}
	} else {
		tips = tip.All(category).Tips
	}

	c.JSON(http.StatusOK, tips)
}

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

func CreateTip(c *gin.Context) {
	form := TipForm{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := user.FindById(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
