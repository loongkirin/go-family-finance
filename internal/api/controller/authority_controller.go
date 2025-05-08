package controller

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loongkirin/gdk/cache"
	"github.com/loongkirin/gdk/cache/redis"
	"github.com/loongkirin/gdk/captcha"
	"github.com/loongkirin/gdk/database/gorm/repository"
	"github.com/loongkirin/gdk/net/http/gin/middleware"
	"github.com/loongkirin/gdk/net/http/request"
	"github.com/loongkirin/gdk/net/http/response"
	"github.com/loongkirin/gdk/oauth"
	"github.com/loongkirin/go-family-finance/internal/app"
	"github.com/loongkirin/go-family-finance/internal/domain/auth"
	"github.com/mojocn/base64Captcha"
)

var (
	cpCache cache.CacheStore
	store   base64Captcha.Store
	cp      *captcha.Captcha

	errCaptchaWrong = errors.New("验证码错误")
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

	if verfied, err := captcha.VerifyCaptcha(store, l.Data.Captcha.CaptchaId, l.Data.Captcha.CaptchaValue, true); err != nil || !verfied {
		if err != nil {
			response.Fail(c, err.Error(), map[string]interface{}{})
		} else {
			response.Fail(c, errCaptchaWrong.Error(), map[string]interface{}{})
		}
		return
	}

	ctx := c.Copy()
	ctx.Set("user_agent", c.Request.UserAgent())
	ctx.Set("client_ip", c.ClientIP())
	r, err := t.authService.Login(ctx, &l)
	if err != nil {
		response.Fail(c, err.Error(), map[string]interface{}{})
		return
	}

	response.Ok(c, "登录成功", r)

}

func (t *AuthorityController) Register(c *gin.Context) {
	var l request.DataRequest[auth.RegisterDTO]

	if err := middleware.ValidateRequest(c, &l); err != nil {
		response.BadRequest(c, err.Error(), map[string]interface{}{})
		return
	}

	// if err := c.ShouldBindJSON(&l); err != nil {
	// 	fmt.Println("2-1", err.Error())
	// 	response.BadRequest(c, "Bad Request:Invalid Parameters", map[string]interface{}{})
	// 	return
	// }

	if verfied, err := captcha.VerifyCaptcha(store, l.Data.Captcha.CaptchaId, l.Data.Captcha.CaptchaValue, true); err != nil || !verfied {
		if err != nil {
			response.Fail(c, err.Error(), map[string]interface{}{})
		} else {
			response.Fail(c, errCaptchaWrong.Error(), map[string]interface{}{})
		}
		return
	}

	r, err := t.authService.Register(c, &l)
	if err != nil {
		response.Fail(c, err.Error(), map[string]interface{}{})
		return
	}

	response.Ok(c, "注册成功", r)
}
