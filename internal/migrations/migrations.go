package migrations

import (
	"github.com/loongkirin/go-family-finance/internal/domain/auth"
	"gorm.io/gorm"
)

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) {
	auth.Migrate(db)
}
