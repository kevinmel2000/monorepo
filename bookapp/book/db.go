package book

// db.go contains all data layer in handling book

const saveBookQuery = `INSERT INTO book(title, author) VALUES($1,$2)`

// saveBook to save book to database
// why separate this to another function? This data layer function should not be mixed with business logic
// data should be treated as it is, the way business logic handle/transform the data should not affect this
func (bs *BookService) saveBook(book Book) error {
	p, err := bs.masterDB.Preparex(saveBookQuery)
	if err != nil {
		return err
	}
	_, err = p.Exec(book.Title, book.Author)
	return err
}
