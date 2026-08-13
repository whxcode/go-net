package model

import (
	"fmt"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func InitDB() {
	// ✅ 正确（用 @tcp）
	dsn := "root:123456@tcp(127.0.0.1:3306)/go-net?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		panic(err)
	}

	DB = db
}
