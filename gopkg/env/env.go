package env

import (
	"bufio"
	"os"
	"runtime"
	"strings"

	"github.com/subosito/gotenv"
	"github.com/lab46/monorepo/gopkg/log"
)

type ServiceEnv string

// Env list
const (
	DevelopmentEnv ServiceEnv = "development"
	StagingEnv     ServiceEnv = "staging"
	ProductionEnv  ServiceEnv = "production"
)

// Env related var
var (
	envName      = "TKPENV"
	currentBuild = "unavailable"
	goVersion    string
)

// env package will read .env file when applicatino is started

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Debug(err)
	}
	goVersion = runtime.Version()
}

// SetFromEnvFile load file
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
		text = strings.TrimSpace(text)
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

func GetEnvName() string {
	return envName
}

// SetEnvName to set env variable name
func SetEnvName(name string) {
	envName = name
}

func EnvList() []ServiceEnv {
	return []ServiceEnv{DevelopmentEnv, StagingEnv, ProductionEnv}
}

func SetCurrentServiceEnv(env ServiceEnv) error {
	return Setenv(envName, string(env))
}

func GetCurrentServiceEnv() string {
	e := Getenv(envName)
	if e == "" {
		e = string(DevelopmentEnv)
	}
	return e
}

func Getenv(key string) string {
	return os.Getenv(key)
}

func Setenv(key, value string) error {
	return os.Setenv(key, value)
}

// SetCurrentBuild to determine the latest build of
func SetCurrentBuild(buildnumber string) {
	currentBuild = buildnumber
}

// GetCurrentBuild return the current build number
func GetCurrentBuild() string {
	return currentBuild
}

// GetGoVersion to return current build go version
func GetGoVersion() string {
	return goVersion
}
