package env

import (
	"bufio"
	"os"
	"strings"

	"github.com/lab46/example/pkg/log"
)

const (
	DevelopmentEnv = "dev"
	StagingEnv     = "staging"
	ProductionEnv  = "prod"
)

// env package will read .env file when applicatino is started

func init() {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return
	} else if err != nil {
		// write something
		log.Debug(err)
		return
	}

	f, err := os.Open(envFile)
	if err != nil {
		log.Debug(err)
		return
	}
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		log.Debug(err)
		return
	}
	for scanner.Scan() {
		text := scanner.Text()
		vars := strings.SplitN(text, "=", 2)
		if len(vars) < 2 {
			return
		}
		if err := Setenv(vars[0], vars[1]); err != nil {
			log.Debug(err)
			return
		}
	}
}

func GetCurrentServiceEnv() string {
	key := "EXMPLENV"
	e := Getenv(key)
	if e == "" {
		e = DevelopmentEnv
	}
	return e
}

func Getenv(key string) string {
	return os.Getenv(key)
}

func Setenv(key, value string) error {
	return os.Setenv(key, value)
}
