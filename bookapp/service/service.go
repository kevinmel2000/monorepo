package service

import (
	"time"

	"github.com/lab46/example/pkg/webserver"
)

type Service struct {
	webserver *webserver.WebServer
}

func New(httpPort string) Service {
	w := webserver.New(webserver.Options{Port: httpPort, Timeout: time.Second * 2})
	service := Service{
		webserver: w,
	}
	return service
}

func (s *Service) RunWebserver() error {
	return s.webserver.Run()
}

func (s *Service) RunGrpcServer() error {
	return nil
}
