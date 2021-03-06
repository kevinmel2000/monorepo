package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/lab46/monorepo/gopkg/env"
	"github.com/lab46/monorepo/gopkg/flags"
	"github.com/lab46/monorepo/gopkg/log"
	"github.com/lab46/monorepo/gopkg/sql/sqldb"
	"github.com/lab46/monorepo/svc/rentapp/service"
)

func initService() (*service.Service, error) {
	serviceConfig, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	_, err = sqldb.Open("postgres", serviceConfig.Postgres.MasterExampleDB)
	if err != nil {
		return nil, err
	}
	_, err = sqldb.Open("postgres", serviceConfig.Postgres.SlaveExampleDB)
	if err != nil {
		return nil, err
	}

	// create new service
	s := service.New("9000")
	return s, err
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
	env.SetConfigDir(sf.configDir)

	s, err := initService()
	if err != nil {
		log.Fatal(err)
	}

	fatalChan := make(chan error)
	go func() {
		fatalChan <- s.RunWebserver()
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
