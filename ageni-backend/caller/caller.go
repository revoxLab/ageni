package caller

import (
	"github.com/readonme/open-studio/conf"
	"github.com/sashabaranov/go-openai"
)

var (
	OpenAIClient *openai.Client
)

func InitCaller(c *conf.Config) {
	OpenAIClient = openai.NewClient(conf.Conf.OpenAIKey)

}
