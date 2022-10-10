package controllers

import (
	"fmt"
	"net/http"

	"dailypractice/user"
	"dailypractice/utils"
	"dailypractice/utils/token"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	RegisterInput
}

func SignUp(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	newUser, ok := user.Save(input.Email, input.Password)

	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "validated!", "user": newUser})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "register failed"})
	}
}

func Login(c *gin.Context) {
	loginForm := LoginInput{}
	c.Bind(&loginForm)

	email := loginForm.Email
	plainPW := loginForm.Password

	currentUser, err := user.FindByEmail(email)
	utils.CheckError(err)

	if currentUser.AuthByPassword(plainPW) {
		jwtToken, err := token.GenerateToken(string(currentUser.Id))

		if err != nil {
			fmt.Printf("error!: %+v", err)
		}
		c.SetCookie("token", jwtToken, 60*60*24, "", "daily-practice.ky2020.click", true, true)
		c.JSON(http.StatusOK, gin.H{
			"message": "Login sucessfully",
			"user":    map[string]string{"email": currentUser.Email},
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "error",
		})
	}
}
