package config

import (
	"flag"

	"github.com/lab46/example/pkg/env"
	"github.com/lab46/example/pkg/flags"
)

type config struct {
	Path string
	Env  string
}

var cfg config

func (c *config) Parse(fs *flag.FlagSet, args []string) error {
	fs.StringVar(&c.Path, "config_path", "", "configuration path")
	return nil
}

func init() {
	flags.Parse(&cfg)
	cfg.Env = env.GetCurrentServiceEnv()
}

func OverrideConfigPath(path string) {
	cfg.Path = path
}

func GetPath() string {
	return cfg.Path
}

func LoadYamlConfig(result interface{}, filename string) error {
	return nil
}
