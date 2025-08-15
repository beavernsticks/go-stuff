package bsgostuff_config

type NATS struct {
	URL         string `env:"INFRASTRUCTURE__NATS__URL"`
	TopicPrefix string `env:"INFRASTRUCTURE__NATS__TOPIC_PREFIX"`
}
