package book

// book.go contains all business logic in handling book

type Book struct {
	ID     int64  `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
}

func AddBook() error {
	return nil
}
