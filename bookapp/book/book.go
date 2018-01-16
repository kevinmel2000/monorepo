package book

// book.go contains all business logic in handling book

import (
	"github.com/lab46/example/pkg/errors"
)

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
func (bs *BookService) AddBook(book Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	// data layer to save a book
	return bs.saveBook(book)
}

func (bs *BookService) ListOfBooks() ([]Book, error) {
	books, err := bs.getBooks()
	return books, err
}
