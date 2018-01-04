package book

import (
	"github.com/jmoiron/sqlx"
)

type BookService struct {
	DB *sqlx.DB
}
