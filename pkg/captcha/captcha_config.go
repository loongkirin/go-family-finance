package captcha

type CaptchaType string

const (
	Audio   CaptchaType = "audio"
	String  CaptchaType = "string"
	Math    CaptchaType = "math"
	Chinese CaptchaType = "chinese"
	Digit   CaptchaType = "digit"
)

type CaptchaConfig struct {
	CaptchaType   CaptchaType `mapstructure:"captcha_type" json:"captcha_type" yaml:"captcha_type"`
	CaptchaLength int         `mapstructure:"captcha_length" json:"captcha_length" yaml:"captcha_length"`
}
