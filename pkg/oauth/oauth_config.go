package oauth

type OAuthConfig struct {
	SecretKey          string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	AccessExpiresTime  string `mapstructure:"access_expires_time" json:"access_expires_time" yaml:"access_expires_time"`
	RefreshExpiresTime string `mapstructure:"refresh_expires_time" json:"refresh_expires_time" yaml:"refresh_expires_time"`
	Issuer             string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
}
