package conf

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type ConfigOption struct {
	sourceFilePath string
}

type ConfigOptionFunc func(opt *ConfigOption)

type ConfigLoader interface {
	Load(data interface{}, opt ...ConfigOptionFunc) error
}

func WithSourceFilePath(source string) ConfigOptionFunc {
	return func(opt *ConfigOption) {
		opt.sourceFilePath = source
	}
}

func NewYamlConfig() *FileConfig {
	return NewFileConfig(func(bytes []byte, i interface{}) error {
		return yaml.Unmarshal(bytes, i)
	})
}

func NewTomlConfig() *FileConfig {
	return NewFileConfig(func(bytes []byte, i interface{}) error {
		_, err := toml.Decode(string(bytes), i)
		return err
	})
}
