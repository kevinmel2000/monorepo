package service

import (
	"time"

	"github.com/lab46/example/gopkg/webserver"
)

type Service struct {
	webserver *webserver.WebServer
}

func New(address string) *Service {
	w := webserver.New(webserver.Options{Address: address, Timeout: time.Second * 2})
	service := Service{
		webserver: w,
	}
	return &service
}

func (s *Service) RunWebserver() error {
	return s.webserver.Run()
}

func (s *Service) RunGrpcServer() error {
	return nil
}
