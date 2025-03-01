package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loongkirin/go-family-finance/internal/app"
	"github.com/loongkirin/go-family-finance/pkg/cache"
	"github.com/loongkirin/go-family-finance/pkg/captcha"
	"github.com/loongkirin/go-family-finance/pkg/http/response"
)

var cpCache = cache.NewInMemoryStore(time.Minute * 3)
var store = captcha.NewCaptchaStore(cpCache, time.Minute*1)
var cp = captcha.NewCaptcha((store))

type AuthorityController struct {
}

func NewAuthorityController() *AuthorityController {
	return &AuthorityController{}
}

func (t *AuthorityController) Captcha(c *gin.Context) {
	if id, b64s, _, err := cp.GenerateDigitCaptcha(app.AppContext.APP_CONFIG.CaptchaConfig.CaptchaLength); err != nil {
		response.Fail(c, "验证码获取失败", map[string]interface{}{})
	} else {
		response.Ok(c, "验证码获取成功", gin.H{
			"captcha_id":     id,
			"pic_path":       b64s,
			"captcha_length": app.AppContext.APP_CONFIG.CaptchaConfig.CaptchaLength,
		})
	}
}
