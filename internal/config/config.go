package config

type AppConfig struct {
	Environment string `mapstructure:"-"`

	AuthConfig AuthConfig `mapstructure:"auth"`
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
