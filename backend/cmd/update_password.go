package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Password string
}

func (User) TableName() string {
	return "sys_user"
}

func main() {
	// 生成密码hash
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成密码hash失败:", err)
	}

	fmt.Println("密码:", password)
	fmt.Println("BCrypt Hash:", string(hash))

	// 连接数据库
	dsn := "warehouse:Warehouse@2026@tcp(localhost:3306)/warehouse?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 更新admin用户密码
	result := db.Model(&User{}).Where("username = ?", "admin").Update("password", string(hash))
	if result.Error != nil {
		log.Fatal("更新密码失败:", result.Error)
	}

	fmt.Printf("\n✅ 成功更新 admin 用户密码！(影响 %d 行)\n", result.RowsAffected)
}
