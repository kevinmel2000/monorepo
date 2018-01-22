package service

func (s *Service) registerHTTPApi() {
	// create subrouter for api/book/v1
	r := s.webserver.Router()
	r = r.SubRouter("/api/v1")

	r.Post("/add", addBook)
	r.Get("/list", bookList)
	r.Get("/book", getBookByID)
}
