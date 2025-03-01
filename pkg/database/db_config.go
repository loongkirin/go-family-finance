package database

type DbConfig struct {
	DbType string         `mapstructure:"db_type" json:"db_type" yaml:"db_type"`
	Master DBConnection   `mapstructure:"master" json:"master" yaml:"master"`
	Slaves []DBConnection `mapstructure:"slaves" json:"slaves" yaml:"slaves"`
}

type DBConnection struct {
	Host            string `mapstructure:"host" json:"host" yaml:"host"`
	Port            int    `mapstructure:"port" json:"port" yaml:"port"`
	Config          string `mapstructure:"config" json:"config" yaml:"config"`
	User            string `mapstructure:"user" json:"user" yaml:"user"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	DbName          string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}
