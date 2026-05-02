package cfg

type HttpCfg struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
