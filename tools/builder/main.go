package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/lab46/monorepo/tools/builder/task"
)

func main() {

	opt, err := task.ReadOptionFromPath("example.yaml")
	if err != nil {
		log.Fatalln(err.Error())
	}

	spew.Dump(opt)
}
