package cfg

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// `AppEnv` is the application runtime environment discriminator.
type AppEnv string

// Supported application runtime environments.
const (
	// `EnvDev` enables development-oriented behavior.
	EnvDev AppEnv = "dev"
	// `EnvProd` enables production-oriented behavior.
	EnvProd AppEnv = "prod"
)

// `GetAppEnv` loads the runtime environment from `APP_ENV`.
func GetAppEnv() AppEnv {
	envRaw := strings.TrimSpace(viper.GetString("APP_ENV"))
	if envRaw == "" {
		return EnvDev
	}

	env := AppEnv(strings.ToLower(envRaw))
	switch env {
	case EnvDev, EnvProd:
		return env
	default:
		panic(fmt.Sprintf("[GetAppEnv] invalid APP_ENV: %s", envRaw))
	}
}
