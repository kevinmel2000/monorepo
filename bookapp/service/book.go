package service

import (
	"net/http"

	"github.com/lab46/example/bookapp/book"
	"github.com/lab46/example/pkg/httpresponse"
)

func addBook(w http.ResponseWriter, r *http.Request) {

}

func bookList(w http.ResponseWriter, r *http.Request) {
	books, err := book.MustGet().ListOfBooks()
	if err != nil {
		httpresponse.InternalServerError(w, err.Error())
	}
	httpresponse.WithData(w, books)
}
