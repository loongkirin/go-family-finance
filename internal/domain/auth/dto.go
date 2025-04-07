package auth

type GeneratedCaptchaDTO struct {
	CaptchaId     string `json:"captcha_id"`
	PicPath       string `json:"pic_path"`
	CaptchaLength int    `json:"captcha_length"`
}

type CaptchaDTO struct {
	CaptchaId    string `json:"captcha_id"`
	CaptchaValue string `json:"captcha_value"`
}

type TenantDTO struct {
	TenantId   string `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
	CreatedAt  int64  `json:"created_at"`
}

type OAuthDTO struct {
	SessionId    string `json:"session_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredAt    int64  `json:"expired_at"`
}

type UserDTO struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantDTO
	OAuthDTO
	CaptchaDTO
}

type CreateUserDTO struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserDTO struct {
	TenantId string `json:"tenant_id" binding:"required"`
	Id       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"omitempty"`
	Phone    string `json:"phone" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty"`
	Active   bool   `json:"active" binding:"omitempty"`
}

type LoginDTO struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CaptchaDTO
}

type RegisterDTO struct {
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	UserName   string `json:"user_name"`
	TenantName string `json:"tenant_name"`
	Password   string `json:"password"`
	CaptchaDTO
}
