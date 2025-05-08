package app

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	validationMessages sync.Map
	validatorOnce      sync.Once
)

var (
	PasswordValidator = NewCustomValidator("password", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		hasNumber := strings.ContainsAny(value, "0123456789")
		hasLetter := strings.ContainsAny(value, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		hasSpecial := strings.ContainsAny(value, "!@#$%^&*()_+-=[]{}|;:,.<>?")
		return len(value) >= 8 && len(value) <= 20 && hasNumber && hasLetter && hasSpecial
	})

	ChineseMobileValidator = NewCustomValidator("cnmobile", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) == 11
	})

	ChineseIdcardValidator = NewCustomValidator("cnidcard", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) == 18
	})

	MinLenValidator = NewCustomValidator("min_len", func(fl validator.FieldLevel) bool {
		param := fl.Param()
		length := len(fl.Field().String())
		minLen, err := strconv.Atoi(param)
		if err != nil {
			return false
		}
		return length >= minLen
	})

	MaxLenValidator = NewCustomValidator("max_len", func(fl validator.FieldLevel) bool {
		param := fl.Param()
		length := len(fl.Field().String())
		maxLen, err := strconv.Atoi(param)
		if err != nil {
			return false
		}
		return length <= maxLen
	})
)

func init() {
	defaultMessages := map[string]string{
		"required": "字段是必需的",
		"min":      "值必须大于或等于%v",
		"max":      "值必须小于或等于%v",
		"len":      "长度必须等于%v",
		"email":    "必须是有效的电子邮件地址%s",
		"mobile":   "必须是有效的手机号%s",
		"idcard":   "必须是有效的身份证号%s",
		"password": "密码必须包含数字和字母，长度在8-20之间%s",
		"min_len":  "长度必须大于或等于%v",
		"max_len":  "长度必须小于或等于%v",
	}

	for k, v := range defaultMessages {
		validationMessages.Store(k, v)
	}
}

// CustomValidator represents a custom validator
type CustomValidator struct {
	Tag          string
	ValidateFunc validator.Func
}

func NewCustomValidator(tag string, validateFunc validator.Func) CustomValidator {
	return CustomValidator{
		Tag:          tag,
		ValidateFunc: validateFunc,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("字段: %s,错误: %s", ve.Field, ve.Message)
}

// ValidationErrors represents a slice of validation errors
type ValidationErrors []ValidationError

func (ves ValidationErrors) Error() string {
	var builder strings.Builder
	for _, ve := range ves {
		builder.WriteString(ve.Error())
		builder.WriteString("\n")
	}
	return builder.String()
}

// RegisterCustomValidators registers custom validators
func RegisterCustomValidators(v *validator.Validate, validators ...CustomValidator) {
	for _, validator := range validators {
		v.RegisterValidation(validator.Tag, validator.ValidateFunc)
	}
}

// RegisterCustomMessages registers custom validation error messages
func RegisterCustomMessages(msg map[string]string) {
	for k, v := range msg {
		validationMessages.Store(k, v)
	}
}

// GetCustomMessage gets the validation error message for a specific tag
func GetCustomMessage(tag string) string {
	if msg, ok := validationMessages.Load(tag); ok {
		return msg.(string)
	}
	return fmt.Sprintf("validate error: %s", tag)
}

// Validator middleware
func Validator(validators ...CustomValidator) gin.HandlerFunc {
	validatorOnce.Do(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			RegisterCustomValidators(v, validators...)
			v.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			})
		}
	})

	return func(c *gin.Context) {
		c.Next()
	}
}

// 验证请求
func ValidateRequest(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errs := make(ValidationErrors, 0, len(validationErrors))

			for _, e := range validationErrors {
				field := e.Field()
				tag := e.Tag()
				value := e.Value()
				param := e.Param()

				fmt.Println("field", field)
				fmt.Println("tag", tag)
				fmt.Println("value", value)
				fmt.Println("param", param)
				fmt.Println("actualTag", e.ActualTag())
				fmt.Println("error", e.Error())

				message := GetCustomMessage(tag)

				errs = append(errs, ValidationError{
					Field:   field,
					Tag:     tag,
					Value:   fmt.Sprintf("%v", value),
					Message: fmt.Sprintf(message, param),
				})
			}

			return errors.New(errs.Error())
		}

		return err
	}
	return nil
}
