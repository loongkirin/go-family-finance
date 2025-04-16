package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserExists         = errors.New("用户已存在")
	ErrUserNotActive      = errors.New("用户未激活")
	ErrTenantExists       = errors.New("租户已存在")
	ErrTenantNotFound     = errors.New("租户不存在")
	ErrTenantNameRequired = errors.New("租户名称不能为空")
	ErrPasswordInvalid    = errors.New("密码错误")
	ErrPhoneRequired      = errors.New("手机号不能为空")
)
