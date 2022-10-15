package token

import (
	"errors"
	"fmt"
	"os"

	"time"

	"dailypractice/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var (
	JWT_SECRET string
)

func init() {
	err := godotenv.Load()
	utils.CheckError(err)
	JWT_SECRET = os.Getenv("JWT_SECRET")
}

func GenerateToken(user_id string) (string, error) {
	token_lifespan := 24
	tokenAndClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
	})

	jwtToken, err := tokenAndClaim.SignedString([]byte(JWT_SECRET))

	return jwtToken, err
}

func TokenValid(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	if tokenString == "" {
		return "", errors.New("Can't find token")
	}
	tokenWithClaim, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if claims, ok := tokenWithClaim.Claims.(jwt.MapClaims); ok && tokenWithClaim.Valid {
		user_id := claims["user_id"].(string)
		return user_id, nil
	} else {
		return "", err
	}
}

func ExtractToken(c *gin.Context) string {
	tokenString := c.Query("token")
	if tokenString != "" {
		return tokenString
	}
	bearToken, err := c.Cookie("token")
	if err != nil {
		utils.CheckError(err)
		return ""
	}

	return bearToken
}

func ExtractTokenID(c *gin.Context) (string, error) {
	user_id, err := TokenValid(c)
	if err != nil {
		return "", err
	}

	return user_id, nil
}
