package main

import (
	"github.com/lab46/monorepo/gopkg/env"
	"github.com/lab46/monorepo/gopkg/redis"
	"github.com/lab46/monorepo/gopkg/sql/sqldb"
)

type ServiceConfig struct {
	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	MasterExampleDB sqldb.Config `yaml:"masterexample"`
	SlaveExampleDB  sqldb.Config `yaml:"slaveexample"`
}

type RedisConfig struct {
	ExampleRedis redis.Config `yaml:"redisexample"`
}

// LoadConfig for loading service configuration
func LoadConfig() (ServiceConfig, error) {
	conf := ServiceConfig{
		Postgres: PostgresConfig{},
		Redis:    RedisConfig{},
	}
	err := env.LoadYamlConfig(&conf.Postgres, "database.yml")
	if err != nil {
		return conf, err
	}
	err = env.LoadYamlConfig(&conf.Redis, "redis.yml")
	if err != nil {
		return conf, err
	}
	return conf, nil
}
