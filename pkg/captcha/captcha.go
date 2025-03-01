package captcha

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	store base64Captcha.Store
}

func NewCaptcha(store base64Captcha.Store) *Captcha {
	return &Captcha{
		store: store,
	}
}

func (c *Captcha) GenerateCaptcha(cfg CaptchaConfig) (id, b64s, answer string, err error) {
	switch cfg.CaptchaType {
	case Audio:
		return c.GenerateAudioCaptcha(cfg.CaptchaLength)
	case String:
		return c.GenerateStringCaptcha(cfg.CaptchaLength)
	case Math:
		return c.GenerateMathCaptcha(cfg.CaptchaLength)
	case Chinese:
		return c.GenerateChinseCaptcha(cfg.CaptchaLength)
	default:
		return c.GenerateDigitCaptcha(cfg.CaptchaLength)
	}
}

func (c *Captcha) GenerateStringCaptcha(length int) (id, b64s, answer string, err error) {
	driver := base64Captcha.NewDriverString(46, 140, 2, 2, length, "234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, nil, []string{"wqy-microhei.ttc"})
	cap := base64Captcha.NewCaptcha(driver, c.store)
	return cap.Generate()
}

func (c *Captcha) GenerateChinseCaptcha(length int) (id, b64s, answer string, err error) {
	driver := base64Captcha.NewDriverChinese(46, 140, 2, 2, length, "234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, nil, []string{"wqy-microhei.ttc"})
	cap := base64Captcha.NewCaptcha(driver, c.store)
	return cap.Generate()
}

func (c *Captcha) GenerateMathCaptcha(length int) (id, b64s, answer string, err error) {
	driver := base64Captcha.NewDriverMath(46, 140, 2, 2, &color.RGBA{240, 240, 246, 246}, nil, []string{"wqy-microhei.ttc"})
	cap := base64Captcha.NewCaptcha(driver, c.store)
	return cap.Generate()
}

func (c *Captcha) GenerateDigitCaptcha(length int) (id, b64s, answer string, err error) {
	driver := base64Captcha.NewDriverDigit(80, 240, length, 0.7, 80)
	cap := base64Captcha.NewCaptcha(driver, c.store)
	return cap.Generate()
}

func (c *Captcha) GenerateAudioCaptcha(length int) (id, b64s, answer string, err error) {
	driver := base64Captcha.NewDriverAudio(length, "zh")
	cap := base64Captcha.NewCaptcha(driver, c.store)
	return cap.Generate()
}

func (c *Captcha) Verify(id, code string, clear bool) bool {
	return c.store.Verify(id, code, clear)
}
