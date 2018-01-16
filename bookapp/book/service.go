package book

import (
	"github.com/jmoiron/sqlx"
	"github.com/lab46/example/pkg/errors"
	"github.com/lab46/example/pkg/log"
	"github.com/lab46/example/pkg/rdbms"
)

var bookService *BookService

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

func Get() (*BookService, error) {
	if bookService == nil {
		return nil, errors.New("Book service is empty, please init the package first")
	}
	return bookService, nil
}

func MustGet() *BookService {
	bs, err := Get()
	if err != nil {
		log.Fatal(err)
	}
	return bs
}
