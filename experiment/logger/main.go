package main

import (
	"github.com/lab46/monorepo/gopkg/log"
)

func main() {
	err := log.SetOutputToFile("experiment.log")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("this is some log")
	log.Println("this is some other log")

	// change log level
	log.SetLevel(log.DebugLevel)
	log.Debug("this is level from debug")
	log.Debug("some more from debug")
}
