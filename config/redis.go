package bsgostuff_config

type Redis struct {
	Host     string `env:"INFRASTRUCTURE__REDIS__HOST"`
	Password string `env:"INFRASTRUCTURE__REDIS__PASSWORD"`
	Database int    `env:"INFRASTRUCTURE__REDIS__DATABASE"`
}
