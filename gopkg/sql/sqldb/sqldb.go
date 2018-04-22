package sqldb

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lab46/monorepo/gopkg/log"
)

// Config of database
type Config struct {
	DSN                string `yaml:"dsn"`
	MaxConnections     int    `yaml:"maxconns"`
	MaxIdleConnections int    `yaml:"maxidleconns"`
	Pretend            bool   `yaml:"pretend"`
	Retry              int    `yaml:"retry"`
}

// Connect to database with default connect timeout for 5 seconds
func Connect(driver, dsn string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, driver, dsn)
	return db, err
}

// Open database connection with sqldb.Config
func Open(driver string, config Config) (*sqlx.DB, error) {
	if config.Pretend {
		db := &sqlx.DB{}
		return db, nil
	}
	log.Debugf("[sqldb][config] %+v", config)

	var (
		err error
		db  *sqlx.DB
	)
	// retry mechanism
	for x := 0; x < config.Retry; x++ {
		db, err = Connect(driver, config.DSN)
		if err == nil {
			break
		} else {
			// log error
			log.Warnf("[sqldb][failed] failed to connect to %s with error %s", config.DSN, err.Error())
		}
		// continue with condition
		log.Warnf("[sqldb][retry] retrying to connect to %s", config.DSN)
		if x+1 == config.Retry && err != nil {
			log.Errorf("[sqldb][error] retry time exhausted, cannot connect to database: %s", err.Error())
			return nil, fmt.Errorf("Failed connect to database: %s", err.Error())
		}
		// sleep for 5 secs everytime retries
		time.Sleep(time.Second * 5)
	}

	// test by pinging database
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if config.MaxConnections > 0 {
		db.SetMaxOpenConns(config.MaxConnections)
	}
	if config.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(config.MaxIdleConnections)
	}
	return db, err
}
