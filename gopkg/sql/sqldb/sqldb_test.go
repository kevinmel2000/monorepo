package sqldb_test

import (
	"testing"

	"github.com/lab46/monorepo/gopkg/sql/sqldb"
)

func TestOpen(t *testing.T) {
	conf := sqldb.Config{
		DSN:     "some_address",
		Pretend: true,
	}
	_, err := sqldb.Open("postgres", conf)
	if err != nil {
		t.Error(err)
	}
}

func TestNewLoadBlanacer(t *testing.T) {
	conf := sqldb.Config{
		DSN:     "some_address",
		Pretend: true,
	}
	db, err := sqldb.Open("postgres", conf)
	if err != nil {
		t.Error(err)
	}
	lb := sqldb.NewLoadBalancer(db)
	lb.GetDB()
}
