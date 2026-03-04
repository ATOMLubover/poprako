package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func Load() (*AppConfig, error) {
	var appConfig AppConfig

	if err := appConfig.load(); err != nil {
		return nil, err
	}

	return &appConfig, nil
}

type AppConfig struct {
	Environment   string `mapstructure:"-"`
	ServerAddress string `mapstructure:"server_address"`

	AuthConfig     AuthConfig     `mapstructure:"auth"`
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
}

func (ac *AppConfig) load() error {
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("app_config")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := v.Unmarshal(&ac); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	environment := os.Getenv("APP_ENVIRONMENT")
	if environment == "" {
		return errors.New("APP_ENVIRONMENT 环境变量未设置")
	}

	ac.Environment = environment

	if err := ac.AuthConfig.load(); err != nil {
		return fmt.Errorf("加载 AuthConfig 失败: %w", err)
	}
	if err := ac.DatabaseConfig.load(); err != nil {
		return fmt.Errorf("加载 DatabaseConfig 失败: %w", err)
	}

	return nil
}

func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

type AuthConfig struct {
	ExpirationHours int    `mapstructure:"expiration_hours"`
	JWTSecretKey    string `mapstructure:"-"`
}

func (ac *AuthConfig) load() error {
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return errors.New("JWT_SECRET_KEY 环境变量未设置")
	}

	ac.JWTSecretKey = jwtSecretKey

	return nil
}

type DatabaseConfig struct {
	DatabaseURL        string `mapstructure:"-"`
	MinIdleConnections int    `mapstructure:"min_idle_connections"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
}

func (dc *DatabaseConfig) load() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL 环境变量未设置")
	}

	dc.DatabaseURL = databaseURL

	return nil
}
