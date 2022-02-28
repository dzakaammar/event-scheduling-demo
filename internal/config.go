package internal

import "github.com/spf13/viper"

type Config struct {
	DbSource    string `mapstructure:"DB_SOURCE"`
	GRPCAddress string `mapstructure:"GRPC_ADDRESS"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
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
