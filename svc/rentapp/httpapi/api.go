package httpapi

import (
	"net/http"

	"github.com/lab46/monorepo/gopkg/router"
)

func RegisterEndpoint(r *router.Router) {
	subBook := r.SubRouter("/rent/v1")
	subBook.Get("/book", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("rent a book"))
	})
}
