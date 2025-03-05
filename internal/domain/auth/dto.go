package auth

type UserDTO struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

type CreateUserDTO struct {
	Name     string `json:"name" binding:"required"`
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserDTO struct {
	TenantId string `json:"tenant_id" binding:"required"`
	Id       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"omitempty"`
	Mobile   string `json:"mobile" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty"`
	Active   bool   `json:"active" binding:"omitempty"`
}
