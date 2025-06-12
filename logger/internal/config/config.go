package config

import "logger/internal/utils"

type MongoDBConfig struct {
	MongoURL string
}

type Config struct {
	WebPort       string
	RPCPort       string
	GRPCPort      string
	MongoDBConfig MongoDBConfig
}

func NewConfig() *Config {
	return &Config{
		WebPort:  utils.String("PORT", "8080"),
		RPCPort:  utils.String("RPC_PORT", "5001"),
		GRPCPort: utils.String("GRPC_PORT", "50001"),
		MongoDBConfig: MongoDBConfig{
			MongoURL: utils.String("MONGODB_URL", "mongodb://mongo:27017"),
		},
	}
}
