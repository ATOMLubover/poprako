package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

// `AppCfg` stores file-based application configuration.
type AppCfg struct {
	Db   *DbCfg   `mapstructure:"database"`
	Http *HttpCfg `mapstructure:"http"`
}

// `NewAppCfg` loads application config from `app_config.json`.
func NewAppCfg() *AppCfg {
	// Load config file.
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("app_config")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("[LoadAppCfg] failed to read config file: %v", err))
	}

	var cfg AppCfg

	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Sprintf("[LoadAppCfg] failed to parse config file: %v", err))
	}

	// Load global environment variables.
	viper.AutomaticEnv()

	return &cfg
}

// `Env` returns the runtime environment from process environment variables.
func (c *AppCfg) Env() AppEnv {
	return GetAppEnv()
}
