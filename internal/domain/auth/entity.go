package auth

import (
	"github.com/loongkirin/go-family-finance/pkg/database"
)

type User struct {
	database.TenantBaseModel
	Name     string `json:"name" gorm:"size:100;not null"`
	Mobile   string `json:"mobile" gorm:"size:100;not null;unique"`
	Password string `json:"password" gorm:"size:100;not null"`
	Active   bool   `json:"active" gorm:"default:true"`
}

func (entity *User) TableName() string {
	return "go_family_finance_user"
}
