package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/lab46/example/pkg/env"
	"github.com/lab46/example/pkg/log"
	"gopkg.in/yaml.v2"
)

type config struct {
	Path string
	Env  string
}

var cfg config

func (c *config) Parse(fs *flag.FlagSet, args []string) error {
	fs.StringVar(&c.Path, "config_path", "", "configuration path")
	return fs.Parse(args)
}

func init() {
	cfg.Env = env.GetCurrentServiceEnv()
}

func SetConfigDir(path string) {
	if f, err := os.Stat(path); err != nil {
		log.Warnf("Failed to check path stats: %s", err.Error())
	} else if !f.IsDir() {
		log.Warnf("%s is not a directory")
	}
	cfg.Path = path
}

func GetPath() string {
	return cfg.Path
}

func LoadYamlConfig(result interface{}, filename string) error {
	dirEnv := strings.ToLower(cfg.Env)
	if dirEnv == "" {
		dirEnv = env.DevelopmentEnv
	}
	confPath := path.Join(cfg.Path, dirEnv, filename)
	log.Debugf("load config from: %s", confPath)
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, result)
}
