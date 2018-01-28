package main

import (
	"testing"

	"github.com/lab46/example/gopkg/env"
)

func TestLoadConfig(t *testing.T) {
	env.SetConfigDir("../files/config/bookapp")
	envList := env.EnvList()
	for _, e := range envList {
		env.SetCurrentServiceEnv(e)
		_, err := LoadConfig()
		if err != nil {
			t.Error(err)
		}
	}
}
