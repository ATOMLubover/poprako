package cfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// AppCfg 是整个应用的根配置
type AppCfg struct {
	Environment   string `json:"-"`
	ServerAddress string `json:"server_address"`

	Auth AuthCfg `json:"auth"`
	DB   DBCfg   `json:"db"`
}

// DBCfg 是数据库的连接配置
type DBCfg struct {
	DSN string `json:"-"`
}

// Load 从 app_config.json 和环境变量中加载应用配置
func Load() (*AppCfg, error) {
	f, err := os.Open("app_config.json")
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %w", err)
	}
	defer f.Close()

	var cfg AppCfg
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	cfg.Environment = env

	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		return nil, errors.New("环境变量 JWT_SECRET_KEY 未设置")
	}
	cfg.Auth.SecretKey = jwtKey

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return nil, errors.New("环境变量 DATABASE_DSN 未设置")
	}
	cfg.DB.DSN = dsn

	return &cfg, nil
}

func (c *AppCfg) IsProduction() bool {
	return c.Environment == "production"
}

func (c *AppCfg) IsDevelopment() bool {
	return c.Environment == "development"
}
