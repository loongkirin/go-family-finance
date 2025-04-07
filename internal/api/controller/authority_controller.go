package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loongkirin/gdk/cache"
	"github.com/loongkirin/gdk/cache/redis"
	"github.com/loongkirin/gdk/captcha"
	"github.com/loongkirin/gdk/database/gorm/repository"
	"github.com/loongkirin/gdk/net/http/request"
	"github.com/loongkirin/gdk/net/http/response"
	"github.com/loongkirin/gdk/oauth"
	"github.com/loongkirin/go-family-finance/internal/app"
	"github.com/loongkirin/go-family-finance/internal/domain/auth"
	"github.com/mojocn/base64Captcha"
)

// var cpCache = cache.NewInMemoryStore(time.Minute * 3)
// var store = captcha.NewCaptchaStore(cpCache, time.Minute*1)
// var cp = captcha.NewCaptcha((store))

var (
	cpCache cache.CacheStore
	store   base64Captcha.Store
	cp      *captcha.Captcha
)

type AuthorityController struct {
	authService auth.AuthService
}

func NewAuthorityController() *AuthorityController {
	cpCache = redis.NewRedisStore(app.AppContext.APP_REDIS.GetMasterDb(), "cpatcha_", time.Minute*3)
	store = captcha.NewCaptchaStore(cpCache, time.Minute*1)
	cp = captcha.NewCaptcha(store)
	oauthMaker, err := oauth.NewPasetoMaker(app.AppContext.APP_CONFIG.OAuthConfig)
	if err != nil {
		panic(err)
	}
	return &AuthorityController{
		authService: auth.NewAuthService(
			repository.NewRepository[auth.User](app.AppContext.APP_DbContext.GetMasterDb()),
			repository.NewRepository[auth.OAuthSession](app.AppContext.APP_DbContext.GetMasterDb()),
			repository.NewRepository[auth.Tenant](app.AppContext.APP_DbContext.GetMasterDb()),
			oauthMaker,
		),
	}
}

func (t *AuthorityController) Captcha(c *gin.Context) {
	if id, b64s, _, err := cp.GenerateDigitCaptcha(app.AppContext.APP_CONFIG.CaptchaConfig.CaptchaLength); err != nil {
		response.Fail(c, "验证码获取失败", map[string]interface{}{})
	} else {
		data := response.DataResponse[auth.GeneratedCaptchaDTO]{
			Data: auth.GeneratedCaptchaDTO{
				CaptchaId:     id,
				PicPath:       b64s,
				CaptchaLength: app.AppContext.APP_CONFIG.CaptchaConfig.CaptchaLength,
			},
		}

		response.Ok(c, "验证码获取成功", data)
	}
}

func (t *AuthorityController) Login(c *gin.Context) {
	var l request.DataRequest[auth.LoginDTO]
	if err := c.ShouldBindJSON(&l); err != nil {
		response.BadRequest(c, "Bad Request:Invalid Parameters", map[string]interface{}{})
		return
	}

	if store.Verify(l.Data.CaptchaId, l.Data.CaptchaValue, true) {
		r, err := t.authService.Login(c, &l)
		if err != nil {
			response.Fail(c, err.Error(), map[string]interface{}{})
			return
		}

		response.Ok(c, "登录成功", r)
	} else {
		response.Fail(c, "验证码错误", map[string]interface{}{})
	}
}

func (t *AuthorityController) Register(c *gin.Context) {
	var l request.DataRequest[auth.RegisterDTO]
	if err := c.ShouldBindJSON(&l); err != nil {
		response.BadRequest(c, "Bad Request:Invalid Parameters", map[string]interface{}{})
		return
	}
	fmt.Println("l.Data.CaptchaId", l.Data.CaptchaId, "l.Data.CaptchaValue", l.Data.CaptchaValue)
	verfied := store.Verify(l.Data.CaptchaId, l.Data.CaptchaValue, true)
	fmt.Println("verfied", verfied)
	if verfied {
		r, err := t.authService.Register(c, &l)
		if err != nil {
			response.Fail(c, err.Error(), map[string]interface{}{})
			return
		}

		response.Ok(c, "注册成功", r)
	} else {
		response.Fail(c, "验证码错误", map[string]interface{}{})
	}
}
