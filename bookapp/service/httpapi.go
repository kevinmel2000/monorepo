package service

import "net/http"

func (s *Service) registerHTTPApi() {
	s.webserver.Router().Get("/something", s.HandlerSomething)
}

func (s *Service) HandlerSomething(w http.ResponseWriter, r *http.Request) {

}
