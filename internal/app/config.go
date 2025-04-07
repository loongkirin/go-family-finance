package app

import (
	"github.com/loongkirin/gdk/cache/redis"
	"github.com/loongkirin/gdk/captcha"
	database "github.com/loongkirin/gdk/database"
	"github.com/loongkirin/gdk/logger"
	"github.com/loongkirin/gdk/oauth"
	"github.com/loongkirin/gdk/telemetry"
)

type ServerConfig struct {
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Mode string `mapstructure:"mode" json:"mode" yaml:"mode"`
}

type AppConfig struct {
	CaptchaConfig   captcha.CaptchaConfig     `mapstructure:"captchaconfig" json:"captchaconfig" yaml:"captchaconfig"`
	OAuthConfig     oauth.OAuthConfig         `mapstructure:"oauthconfig" json:"oauthconfig" yaml:"oauthconfig"`
	RedisConfig     redis.RedisConfig         `mapstructure:"redisconfig" json:"redisconfig" yaml:"redisconfig"`
	DbConfig        database.DbConfig         `mapstructure:"dbconfig" json:"dbconfig" yaml:"dbconfig"`
	ServerConfig    ServerConfig              `mapstructure:"serverconfig" json:"serverconfig" yaml:"serverconfig"`
	LoggerConfig    logger.LoggerConfig       `mapstructure:"loggerconfig" json:"loggerconfig" yaml:"loggerconfig"`
	TelemetryConfig telemetry.TelemetryConfig `mapstructure:"telemetryconfig" json:"telemetryconfig" yaml:"telemetryconfig"`
}
