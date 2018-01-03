package httpapi

import (
	"net/http"

	"github.com/lab46/example/pkg/router"
)

func RegisterEndpoint(r *router.Router) {
	subBook := r.SubRouter("/book/v1")
	subBook.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("list of book"))
	})
	subBook.Get("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("a book"))
	})
}
