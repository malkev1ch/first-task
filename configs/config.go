package configs

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/viper"
	"time"
)

type (
	Config struct {
		Postgres    PostgresConfig
		HTTP        HTTPConfig
		Image       ImageStorage
	}

	PostgresConfig struct {
		Port     string
		Sslmode  string
		Host     string
		Username string
		Dbname   string
		Password string `env:"DB_PASSWORD,unset"`
	}
	HTTPConfig struct {
		Host               string
		Port               string
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxHeaderMegabytes int
	}
	ImageStorage struct {
		PhotoPathCat string
	}
)

func Init(configsDir string) (*Config, error) {

	if err := parseConfigFile(configsDir); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := setFromEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func unmarshal(cfg *Config) error {

	if err := viper.UnmarshalKey("db", &cfg.Postgres); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("image", &cfg.Image); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	fmt.Printf("%+v\n", cfg)
	return nil
}
