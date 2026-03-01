package config

type AuthConfig struct {
	ExpirationHours int    `mapstructure:"expiration_hours"`
	JWTSecretKey    string `mapstructure:"-"`
}
