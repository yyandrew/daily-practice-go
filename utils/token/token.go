package token

import (
	"fmt"
	"os"

	// "os"
	"time"

	"dailypractice/utils"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var (
	JWT_SECRET string
)

func init() {
	err := godotenv.Load()
	utils.CheckError(err)
	JWT_SECRET = os.Getenv("JWT_SECRET")
	fmt.Printf("jwt secret: %s\n", JWT_SECRET)
}

func GenerateToken(user_id string) (string, error) {
	fmt.Printf("user_id: %d, JWT_SECRET: %s\n", user_id, JWT_SECRET)
	token_lifespan := 24
	claims := jwtv4.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	jwtToken := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(JWT_SECRET))
}
