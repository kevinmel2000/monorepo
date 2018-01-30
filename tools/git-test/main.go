package main

import (
	"github.com/lab46/example/gopkg/print"
)

func main() {
	rootCMD := initCMD()
	registerGitTestCommand(rootCMD)
	print.Error(rootCMD.Execute())
}
