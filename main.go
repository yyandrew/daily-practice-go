package main

import (
	"dailypractice/tip"
	"dailypractice/user"
	"dailypractice/utils"
	"encoding/json"
	"fmt"
)

func main() {
	tips := tip.All()
	users := user.All()
	j, err := json.Marshal(users)
	utils.CheckError(err)
	fmt.Printf("users: %s\n", j)

	tj, err := json.Marshal(tips)
	utils.CheckError(err)
	fmt.Printf("json tips: %s\n", tj)
}
