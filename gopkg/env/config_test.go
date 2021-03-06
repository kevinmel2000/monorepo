package env_test

import (
	"testing"

	"github.com/lab46/monorepo/gopkg/env"
)

func TestSetAndGetConfigDir(t *testing.T) {
	dir := "../../files/testfile/yamlconfig"
	err := env.SetConfigDir(dir)
	if err != nil {
		t.Error(err)
	}
	confdir := env.GetConfigDir()
	if dir != confdir {
		t.Errorf("Expecting %s but got %s", dir, confdir)
	}
}

func TestLoadYamlConfig(t *testing.T) {
	configDir := "../../files/testfile/yamlconfig"
	err := env.SetConfigDir(configDir)
	if err != nil {
		t.Error(err)
		return
	}
	tc := struct {
		test struct {
			key1 string `yaml:"key1"`
		}
	}{}
	err = env.LoadYamlConfig(&tc, "test.yaml")
	if err != nil {
		t.Error(err)
	}
}
