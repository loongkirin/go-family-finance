package auth

import (
	"github.com/loongkirin/go-family-finance/internal/domain/models"
)

type User struct {
	models.TenantBaseModel
	Name     string `json:"name" gorm:"size:100;not null"`
	Email    string `json:"email" gorm:"size:100;not null"`
	Phone    string `json:"phone" gorm:"size:100;not null"`
	Password string `json:"password" gorm:"size:100;not null"`
	Active   bool   `json:"active" gorm:"default:true"`
}

func (entity *User) TableName() string {
	return "finance_user"
}

type OAuthSession struct {
	models.TenantBaseModel
	UserId       string `json:"user_id" gorm:"size:32"`
	Email        string `json:"email" gorm:"size:100"`
	Phone        string `json:"phone" gorm:"size:100"`
	UserName     string `json:"user_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	ClientIp     string `json:"client_ip"`
	IsBlocked    bool   `json:"is_blocked"`
	ExpiredAt    int64  `json:"expired_at"`
}

func (entity *OAuthSession) TableName() string {
	return "finance_oauth_session"
}

type Tenant struct {
	models.DbBaseModel
	Name string `json:"name" gorm:"size:500;not null"`
}

func (entity *Tenant) TableName() string {
	return "finance_tenant"
}
