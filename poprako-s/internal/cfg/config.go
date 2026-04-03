package cfg

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// AppCfg 是整个应用的根配置
type AppCfg struct {
	Env      string `mapstructure:"-"`
	HTTPAddr string `mapstructure:"http_address"`

	Auth AuthCfg `mapstructure:"auth"`
	DB   DBCfg   `mapstructure:"db"`
}

// Load 从 app_config.json 和环境变量中加载应用配置
func Load() (*AppCfg, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("app_config")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg AppCfg
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	cfg.Env = env

	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		return nil, errors.New("环境变量 JWT_SECRET_KEY 未设置")
	}
	cfg.Auth.SecretKey = jwtKey

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("环境变量 DATABASE_URL 未设置")
	}
	cfg.DB.DSN = dsn

	return &cfg, nil
}

func (c *AppCfg) IsProduction() bool {
	return c.Env == "production"
}

func (c *AppCfg) IsDevelopment() bool {
	return c.Env == "development"
}
