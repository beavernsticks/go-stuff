package bsgostuff_config

type Redis struct {
	Addr     string `env:"INFRASTRUCTURE__REDIS__ADDR"`
	Password string `env:"INFRASTRUCTURE__REDIS__PASSWORD"`
	Database int    `env:"INFRASTRUCTURE__REDIS__DATABASE"`
	Prefix   string `env:"INFRASTRUCTURE__REDIS__PREFIX"`
}
