package book

// book.go contains all business logic in handling book

import (
	"github.com/jmoiron/sqlx"
	"github.com/lab46/example/gopkg/errors"
	"github.com/lab46/example/gopkg/sqldb"
)

// a package level variable to hold dependencies in the package
var bookService *BookService

type BookService struct {
	masterDB *sqlx.DB
	slaveDB  *sqldb.LoadBalancer
}

func Init(masterDB *sqlx.DB, slave *sqldb.LoadBalancer) *BookService {
	bookService = &BookService{
		masterDB: masterDB,
		slaveDB:  slave,
	}
	return bookService
}

type Book struct {
	ID     int64  `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
}

func (b *Book) Validate() error {
	// check title
	if b.Title == "" {
		return errors.New("book title cannot be empty")
	}
	if len(b.Title) > 30 {
		return errors.New("book title is too long")
	}
	// check author
	if b.Author == "" {
		return errors.New("book author cannot be empty")
	}
	if len(b.Author) > 30 {
		return errors.New("book title is too long")
	}
	return nil
}

// AddBook is a business logic layer to add a book
func AddBook(book Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	// data layer to save a book
	return bookService.saveBook(book)
}

func ListOfBooks() ([]Book, error) {
	books, err := bookService.getBooks()
	return books, err
}

func GetBookByID(id int64) (Book, error) {
	book := Book{}
	book, err := bookService.getBookByID(id)
	return book, err
}
