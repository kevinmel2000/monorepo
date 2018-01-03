package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lab46/example/pkg/webserver"
	"github.com/lab46/example/rentapp/httpapi"
)

func main() {
	w := webserver.New(webserver.Options{})
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
