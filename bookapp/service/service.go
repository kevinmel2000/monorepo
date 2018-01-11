package service

import (
	"time"

	"github.com/lab46/example/bookapp/book"
	"github.com/lab46/example/pkg/webserver"
)

type Service struct {
	webserver    *webserver.WebServer
	dependencies ServiceDependencies
}

type httpAPI struct {
	dependencies *ServiceDependencies
}

type grpcService struct {
	dependencies *ServiceDependencies
}

type ServiceDependencies struct {
	book *book.BookService
}

func New(httpPort string, dependencies ServiceDependencies) Service {
	w := webserver.New(webserver.Options{Port: httpPort, Timeout: time.Second * 2})
	service := Service{
		webserver:    w,
		dependencies: dependencies,
	}
	return service
}

func (s *Service) RunWebserver() error {
	return s.webserver.Run()
}

func (s *Service) RunGrpcServer() error {
	return nil
}
