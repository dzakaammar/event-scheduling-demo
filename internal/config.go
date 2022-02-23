package internal

import "github.com/spf13/viper"

type Config struct {
	DbDriver    string `mapstructure:"DB_DRIVER"`
	DbSource    string `mapstructure:"DB_SOURCE"`
	GRPCAddress string `mapstructure:"GRPC_ADDRESS"`
}

func LoadConfig(path string) (Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
