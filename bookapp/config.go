package main

import (
	"github.com/lab46/example/pkg/config"
	"github.com/lab46/example/pkg/rdbms"
	"github.com/lab46/example/pkg/redis"
)

type ServiceConfig struct {
	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	MasterExampleDB rdbms.Config `yaml:"masterexample"`
	SlaveExampleDB  rdbms.Config `yaml:"slaveexample"`
}

type RedisConfig struct {
	ExampleRedis redis.Config `yaml:"redisexample"`
}

func LoadConfig() (ServiceConfig, error) {
	conf := ServiceConfig{
		Postgres: PostgresConfig{},
		Redis:    RedisConfig{},
	}
	err := config.LoadYamlConfig(&conf.Postgres, "database.yml")
	if err != nil {
		return conf, err
	}
	err = config.LoadYamlConfig(&conf.Redis, "redis.yml")
	if err != nil {
		return conf, err
	}
	return conf, nil
}
