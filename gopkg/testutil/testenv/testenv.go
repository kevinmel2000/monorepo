package testenv

import (
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lab46/monorepo/gopkg/env"
	"github.com/lab46/monorepo/gopkg/log"
)

// EnvConfig to store testing environment configuration
var EnvConfig struct {
	PostgresDSN      string `envconfig:"POSTGRES_DSN" default:"postgres://logistic:logistic@localhost:5432?sslmode=disable"`
	MySQLDSN         string `envconfig:"MYSQL_DSN"`
	RedisAddress     string `envconfig:"REDIS_ADDRESS"`
	MongoDBAddress   string `envconfig:"MONGO_ADDRESS"`
	CassandraAddress string `envconfig:"CASSANDRA_ADDRESS"`
}

func init() {
	err := envconfig.Process("", &EnvConfig)
	if err != nil {
		log.Errorf("Failed to load test env: %s", err.Error())
	}
}
