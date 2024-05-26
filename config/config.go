package config

import (
	"flag"

	"flag"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port    int
		Storage string
	}
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.App.Storage, "storage", "inmemory", "Выбор типа хранилища")
	flag.Parse()
}
