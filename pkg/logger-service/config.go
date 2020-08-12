package logger_service

type Config struct {
	IsVerbose bool `mapstructure:"LOG_VERBOSE"`
}

func (c Config) Validate() error {
	return nil
}
