package auth

import (
	"fmt"

	"gorm.io/gorm"
)

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) {
	// 创建用户表
	if err := db.AutoMigrate(&User{}); err != nil {
		fmt.Println("创建用户表失败", err)
	}

	// 创建OAuthSession表
	if err := db.AutoMigrate(&OAuthSession{}); err != nil {
		fmt.Println("创建OAuthSession表失败", err)
	}

	// 创建Tenant表
	if err := db.AutoMigrate(&Tenant{}); err != nil {
		fmt.Println("创建Tenant表失败", err)
	}

	fmt.Println("Auth模块迁移完成")
}
