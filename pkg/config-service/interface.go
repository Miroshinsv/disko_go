package config_service

const (
	DefaultConfigPath = "configs/config.toml"
)

type IConfig interface {
	Load(path string) error
	Convert(subject ISubjectInterface) error
}
