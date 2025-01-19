package conf

import (
	"os"
)

type FileConfig struct {
	data    interface{}
	opt     ConfigOption
	decoder DecoderFunc
}

type DecoderFunc func([]byte, interface{}) error

func (c *FileConfig) Load(i interface{}, optFunc ...ConfigOptionFunc) error {
	c.data = i
	for _, ofc := range optFunc {
		ofc(&c.opt)
	}
	return c.doLoad()
}

func (c *FileConfig) doLoad() error {
	var (
		err     error
		content []byte
	)

	content, err = os.ReadFile(c.opt.sourceFilePath)

	if err != nil {
		return err
	}

	if err = c.decoder(content, c.data); err != nil {
		return err
	}

	return nil
}

func NewFileConfig(decoder DecoderFunc) *FileConfig {
	return &FileConfig{
		decoder: decoder,
	}
}
