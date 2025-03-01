package redis

type RedisConfig struct {
	Master   RedisConnection   `mapstructure:"master" json:"master" yaml:"master"`
	Slaves   []RedisConnection `mapstructure:"slaves" json:"slaves" yaml:"slaves"`
	PoolSize int               `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
}

type RedisConnection struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}
