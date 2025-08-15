package bsgostuff_config

type PostgreSQL struct {
	Host     string `env:"INFRASTRUCTURE__POSTGRESQL__HOST"`
	Port     string `env:"INFRASTRUCTURE__POSTGRESQL__PORT"`
	User     string `env:"INFRASTRUCTURE__POSTGRESQL__USER"`
	Password string `env:"INFRASTRUCTURE__POSTGRESQL__PASSWORD"`
	DBName   string `env:"INFRASTRUCTURE__POSTGRESQL__DBNAME"`
}
