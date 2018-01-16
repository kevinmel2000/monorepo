package service

func (s *Service) registerHTTPApi() {
	// create subrouter for api/book/v1
	r := s.webserver.Router()
	r = r.SubRouter("/api/book/v1")

	r.Get("/list", bookList)
}
