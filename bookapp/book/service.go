package book

import (
	"github.com/jmoiron/sqlx"
	"github.com/lab46/example/pkg/rdbms"
)

type BookService struct {
	masterDB *sqlx.DB
	slaveDB  *rdbms.LoadBalancer
}

func Init(masterDB *sqlx.DB, slave *rdbms.LoadBalancer) *BookService {
	service := &BookService{
		masterDB: masterDB,
		slaveDB:  slave,
	}
	return service
}

func Handler() {

}
