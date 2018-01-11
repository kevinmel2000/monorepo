package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lab46/example/bookapp/book"
	"github.com/lab46/example/pkg/config"
	"github.com/lab46/example/pkg/flags"
	"github.com/lab46/example/pkg/log"
	"github.com/lab46/example/pkg/rdbms"
	"github.com/lab46/example/pkg/webserver"
	"github.com/lab46/example/rentapp/httpapi"
)

func initDependencies() error {
	serviceConfig, err := LoadConfig()
	if err != nil {
		return err
	}
	masterDB, err := rdbms.Open("postgres", serviceConfig.Postgres.MasterExampleDB)
	if err != nil {
		return err
	}
	slaveDB, err := rdbms.Open("postgres", serviceConfig.Postgres.SlaveExampleDB)
	if err != nil {
		return err
	}
	book.Init(masterDB, rdbms.NewLoadBalancer(slaveDB))
	return err
}

type serviceFlags struct {
	logLevel  string
	configDir string
}

var sf serviceFlags

func (sf *serviceFlags) Parse(fs *flag.FlagSet, args []string) error {
	fs.StringVar(&sf.logLevel, "log_level", "", "logging level")
	fs.StringVar(&sf.configDir, "config_dir", "", "configuration directory")
	return fs.Parse(args)
}

func main() {
	sf = serviceFlags{}
	flags.Parse(&sf)
	log.SetLevelString(sf.logLevel)
	config.SetConfigDir(sf.configDir)

	if err := initDependencies(); err != nil {
		log.Fatal(err)
	}

	w := webserver.New(webserver.Options{
		Port:    "9000",
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
