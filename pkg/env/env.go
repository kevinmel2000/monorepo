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
	err := SetFromEnvFile(".env")
	if err != nil {
		log.Debug(err)
	}
}

func SetFromEnvFile(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return err
	} else if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		return err
	}
	for scanner.Scan() {
		text := scanner.Text()
		vars := strings.SplitN(text, "=", 2)
		if len(vars) < 2 {
			return err
		}
		if err := Setenv(vars[0], vars[1]); err != nil {
			return err
		}
	}
	return nil
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
