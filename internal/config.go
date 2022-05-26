package internal

import "github.com/spf13/viper"

type Config struct {
	DbSource           string `mapstructure:"db_source"`
	GRPCAddress        string `mapstructure:"grpc_address"`
	GRPCGatewayAddress string `mapstructure:"grpc_gateway_address"`
	OTLPEndpoint       string `mapstructure:"otel_exporter_otlp_endpoint"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, nil
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
