package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUri               string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri            string `mapstructure:"REDIS_URL"`
	Port                string `mapstructure:"PORT"`
	JWTSecretKey        string `mapstructure:"JWT_SECRET_KEY"`
	JWTRefreshSecretKey string `mapstructure:"JWT_REFRESH_SECRET_KEY"`
}

var (
	AppConfig *Config
)

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
