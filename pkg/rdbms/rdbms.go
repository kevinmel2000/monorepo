package rdbms

import (
	"github.com/jmoiron/sqlx"
	"github.com/lab46/example/pkg/log"
	_ "github.com/lib/pq"
)

type dsn struct {
	Master string
	Slave  string
}

// Config of database
type Config struct {
	DSN      map[string]dsn
	Skipinit bool
}

// DB of database
type db struct {
	connectedDbs map[dbType]*sqlx.DB
}

var dbObject *db

// dbType is type of database
type dbType string

// Init database connection
func Init(cfg Config) error {
	if cfg.Skipinit {
		return nil
	}

	dbObject = &db{connectedDbs: make(map[dbType]*sqlx.DB)}
	for dbName, dsn := range cfg.DSN {
		log.Debugf("[Database] Connecting to database [%s]...", dbName)
		newDB, err := sqlx.Open("postgres", dsn.Master+";"+dsn.Slave)
		if err != nil {
			log.Errorf("[Database] Failed to connect to db %s. Error: %s", dbName, err.Error())
			return err
		}
		dbObject.connectedDbs[dbType(dbName)] = newDB
	}
	return nil
}
