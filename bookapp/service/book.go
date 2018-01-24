package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lab46/example/bookapp/book"
	"github.com/lab46/example/pkg/http/httpresponse"
)

func addBook(w http.ResponseWriter, r *http.Request) {
	jsonContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpresponse.InternalServerError(w, err.Error())
		return
	}
	b := book.Book{}
	err = json.Unmarshal(jsonContent, &b)
	if err != nil {
		httpresponse.InternalServerError(w, err.Error())
		return
	}
	err = book.AddBook(b)
	if err != nil {
		httpresponse.InternalServerError(w, err.Error())
		return
	}
	httpresponse.StatusOK(w)
}

func getBookByID(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	bookID, err := strconv.ParseInt(id, 64, 0)
	if err != nil {
		httpresponse.InternalServerError(w, err.Error())
		return
	}
	book, err := book.GetBookByID(bookID)
	if err != nil {
		httpresponse.InternalServerError(w, err.Error())
		return
	}
	httpresponse.WithData(w, book)
}

func bookList(w http.ResponseWriter, r *http.Request) {
	books, err := book.ListOfBooks()
	if err != nil {
		httpresponse.InternalServerError(w, err.Error())
		return
	}
	httpresponse.WithData(w, books)
}
