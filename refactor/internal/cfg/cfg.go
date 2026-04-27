package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppCfg struct {
	Env AppEnv

	Db   *DbCfg
	Http *HttpCfg
}

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
