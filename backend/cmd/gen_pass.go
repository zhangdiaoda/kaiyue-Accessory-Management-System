package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	// 使用与后端一致的 DefaultCost (10)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println("Password: " + password)
	fmt.Println("Hash:     " + string(bytes))
}
