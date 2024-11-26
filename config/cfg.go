package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	ExternalAPI string `mapstructure:"HTTP_EXTERNALAPI"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBName      string `mapstructure:"DB_NAME"`
	DBUsername  string `mapstructure:"DB_USERNAME"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
}

func ReadConfig(name, path string) (*Config, error) {

	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
