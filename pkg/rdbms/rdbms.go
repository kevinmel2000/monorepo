package rdbms

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config of database
type Config struct {
	DSN                string `yaml:"dsn"`
	MaxConnections     int    `yaml:"maxconns"`
	MaxIdleConnections int    `yaml:"maxidleconns"`
	Pretend            bool   `yaml:"pretend"`
	Retry              int    `yaml:"retry"`
}

type LoadBalancer struct {
	dbs    []*sqlx.DB
	length int
	count  uint64
}

func Open(driver string, config Config) (*sqlx.DB, error) {
	if config.Pretend {
		db := &sqlx.DB{}
		return db, nil
	}

	var (
		err error
		db  *sqlx.DB
	)
	// retry mechanism
	for x := 0; x < config.Retry; x++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		db, err = sqlx.ConnectContext(ctx, driver, config.DSN)
		if err == nil {
			cancel()
			break
		}
		// else continue with condition
		cancel()
		if x+1 == config.Retry && err != nil {
			return nil, fmt.Errorf("Failed connect to database: %s", err.Error())
		}
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

func NewLoadBalancer(sqlxDbs ...*sqlx.DB) *LoadBalancer {
	l := &LoadBalancer{
		dbs:    sqlxDbs,
		length: len(sqlxDbs),
	}
	return l
}

func (l *LoadBalancer) GetDB() *sqlx.DB {
	return l.dbs[l.get()]
}

// get will return number in db length with round-robin functionality
func (l *LoadBalancer) get() int {
	if l.length <= 1 {
		return 0
	}
	db := int(1 + (atomic.AddUint64(&l.count, 1) % uint64(l.length-1)))
	return db
}
