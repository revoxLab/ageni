package conf

import (
	"flag"
	"github.com/readonme/open-studio/common/auth/jwt"
	"github.com/readonme/open-studio/common/conf"
	"github.com/readonme/open-studio/common/gorm"
	"github.com/readonme/open-studio/common/httpserver"
	"github.com/readonme/open-studio/common/log"
	"math/rand"
	"os"
	"time"
)

var (
	Conf     = &Config{}
	confPath string
)

type Config struct {
	Env        string             `toml:"env"`
	Log        *log.Config        `toml:"log"`
	HttpServer *httpserver.Config `toml:"http_server"`
	StudioDB   *gorm.Config       `toml:"studio_db"`
	JWTToken   *jwt.Config        `toml:"jwt_token"`
	OpenAIKey  string             `toml:"openai_key"`
	BotTabs    []string           `toml:"bot_tabs"`
	PluginTabs []string           `toml:"plugin_tabs"`
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Init config from config file.
func Init() {
	deployEnv := os.Getenv("DEPLOY_ENV")
	log.Infof("deployEnv =%s", deployEnv)
	if deployEnv == "" {
		deployEnv = "local"
	}
	switch deployEnv {
	case "local":
		confPath = "./conf/local.toml"
	case "develop":
		confPath = "./conf/dev.toml"
	case "qa":
		confPath = "./conf/qa.toml"
	case "test":
		confPath = "./conf/test.toml"
	case "staging":
		confPath = "./conf/staging.toml"
	case "prod":
		confPath = "./conf/prod.toml"
	}

	if err := conf.NewTomlConfig().Load(Conf, conf.WithSourceFilePath(confPath)); err != nil {
		panic(err)
	}
	Conf.Env = deployEnv
}
