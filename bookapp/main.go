package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lab46/example/bookapp/httpapi"
	"github.com/lab46/example/pkg/webserver"
)

func main() {
	w := webserver.New(webserver.Options{
		Port:    9000,
		Timeout: time.Second * 2,
	})
	httpapi.RegisterEndpoint(w.Router())

	fatalChan := make(chan error)
	go func() {
		fatalChan <- w.Run()
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Println("Signal terminate detected")
	case err := <-fatalChan:
		log.Fatal("Application failed to run because ", err.Error())
	}
}
