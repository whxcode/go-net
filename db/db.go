package db

import (
	"fmt"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

type panicPlugin struct{}

func (p *panicPlugin) Name() string {
	return "panicPlugin"
}

/**
* 捕获 sql 错误
*
* */
func panicIfError(db *gorm.DB) {
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		fmt.Println("Database operation failed:", db.Error)
		panic(db.Error)
	}
}

func (p *panicPlugin) Initialize(db *gorm.DB) error {
	db.Callback().Query().
		After("gorm:query").
		Register("panic:query", panicIfError)

	db.Callback().Create().
		After("gorm:create").
		Register("panic:create", panicIfError)

	db.Callback().Update().
		After("gorm:update").
		Register("panic:update", panicIfError)

	db.Callback().Delete().
		After("gorm:delete").
		Register("panic:delete", panicIfError)

	db.Callback().Row().
		After("gorm:row").
		Register("panic:row", panicIfError)

	return nil
}

var DB *gorm.DB

func InitDB() {
	// ✅ 正确（用 @tcp）
	dsn := "root:123456@tcp(127.0.0.1:3306)/go-net?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		panic(err)
	}

	db.Use(&panicPlugin{})
	DB = db.Debug()
}

func Must(reuslt *gorm.DB) *gorm.DB {
	if reuslt.Error != nil {
		panic(reuslt.Error)
	}

	return reuslt
}
