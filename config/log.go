package bsgostuff_config

type Log struct {
	Level string `env:"LOG__LEVEL" env-default:"info"`
	Mode  string `env:"LOG__MODE" env-default:"production"`
}
