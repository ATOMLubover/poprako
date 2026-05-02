package cfg

type DbCfg struct {
	// Fields below are automatically loaded from config file.
	MaxOpen     int `mapstructure:"max_open"`
	MaxIdle     int `mapstructure:"max_idle"`
	MaxLife     int `mapstructure:"max_life"`      // unit: second
	MaxIdleTime int `mapstructure:"max_idle_time"` // unit: second
}
