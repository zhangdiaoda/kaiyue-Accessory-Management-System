package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// CheckPassword 验证密码 (复制自 pkg/utils/jwt.go)
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	password := "admin123"
	// 这是 01-schema.sql 和我们刚才手动插入的哈希值
	hash := "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH"

	fmt.Printf("Testing password: %s\n", password)
	fmt.Printf("Against hash: %s\n", hash)

	if CheckPassword(password, hash) {
		fmt.Println("✅ Success: Password matches hash.")
	} else {
		fmt.Println("❌ Failure: Password does NOT match hash.")
		// 尝试生成正确的哈希
		bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		fmt.Printf("ℹ️ Correct hash should be something like: %s\n", string(bytes))
	}
}
