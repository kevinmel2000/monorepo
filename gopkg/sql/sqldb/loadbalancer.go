package sqldb

import (
	"sync/atomic"

	"github.com/jmoiron/sqlx"
)

// LoadBalancer struct
type LoadBalancer struct {
	dbs    []*sqlx.DB
	length int
	count  uint64
}

// NewLoadBalancer create new loadbalancer for database connection
func NewLoadBalancer(sqlxDbs ...*sqlx.DB) *LoadBalancer {
	l := &LoadBalancer{
		dbs:    sqlxDbs,
		length: len(sqlxDbs),
	}
	return l
}

// GetDB return db from loadbalancer
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
