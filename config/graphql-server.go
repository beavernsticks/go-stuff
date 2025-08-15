package bsgostuff_config

type GraphQLServer struct {
	Address              string `env:"GRAPHQL_SERVER__ADDRESS" env-default:":8000"`
	IntrospectionEnabled bool   `env:"GRAPHQL_SERVER__INTROSPECTION_ENABLED" env-default:"false"`
}
