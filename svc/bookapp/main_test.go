package main

import (
	"testing"

	"github.com/lab46/monorepo/gopkg/env"
)

func TestLoadConfig(t *testing.T) {
	err := env.SetConfigDir("files/config")
	if err != nil {
		t.Error(err)
		return
	}
	envList := env.EnvList()
	for _, e := range envList {
		env.SetCurrentServiceEnv(e)
		_, err := LoadConfig()
		if err != nil {
			t.Error(err)
			return
		}
	}
}
