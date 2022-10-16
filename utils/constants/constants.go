package constants

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	IMG_PATH, DOMAIN, JWT_SECRET string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	IMG_PATH = os.Getenv("IMG_PATH")
	DOMAIN = os.Getenv("DOMAIN")
	JWT_SECRET = os.Getenv("JWT_SECRET")
}
