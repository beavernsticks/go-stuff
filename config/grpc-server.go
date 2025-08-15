package bsgostuff_config

type GRPCServer struct {
	Address           string `env:"GRPC_SERVER__ADDRESS" env-default:":50051"`
	ReflectionEnabled bool   `env:"GRPC_SERVER__REFLECTION_ENABLED" env-default:"false"`
}
