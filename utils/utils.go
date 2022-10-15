// Package utils provides ...
package utils

import "fmt"

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("error!: %+v", err)

		return
	}
}
