package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port    int
		Storage string `mapstructure:"STORAGE_TYPE"`
	}
	DB struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     int    `mapstructure:"DB_PORT"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Name     string `mapstructure:"DB_NAME"`
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

	config.getEnv()

	return &config, nil
}

func (c *Config) ParseFlags() {
	stype := os.Getenv("STORAGE_TYPE")
	if stype != "" {
		c.App.Storage = stype
		return
	}

	flag.StringVar(&c.App.Storage, "storage", "inmemory", "Выбор типа хранилища")
	flag.Parse()
}

func (c *Config) getEnv() {
	dbh := os.Getenv("DB_HOST")
	if dbh != "" {
		c.DB.Host = dbh
	}

	dbp := os.Getenv("DB_PORT")
	if dbp != "" {
		var err error
		c.DB.Port, err = strconv.Atoi(dbp)
		if err != nil {
			c.DB.Port = 5432
		}
	}

	dbu := os.Getenv("DB_USER")
	if dbu != "" {
		c.DB.User = dbu
	}

	dbpass := os.Getenv("DB_PASSWORD")
	if dbpass != "" {
		c.DB.Password = dbpass
	}

	dbn := os.Getenv("DB_NAME")
	if dbn != "" {
		c.DB.Name = dbn
	}
}
