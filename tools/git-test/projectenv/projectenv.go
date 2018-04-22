package projectenv

import (
	"os"
	"path"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lab46/monorepo/gopkg/env"
	"github.com/lab46/monorepo/gopkg/log"
)

// default var
var (
	defaultRepoDir = path.Join(os.Getenv("GOPATH"), "src", "github.com", "lab46", "monorepo")
)

// Config to store project environment configuration
var Config struct {
	RepoName      string `envconfig:"GT_REPO_NAME" default:"lab46/monorepo"`
	RepoDir       string `envconfig:"GT_REPO_DIR"`
	ServiceFolder string `envconfig:"GT_SERVICE_FOLDER" default:"svc"`
	Env           string `envconfig:"GT_ENV" default:""`
}

func replaceIfEmpty(s *string, replacer string) {
	if *s == "" {
		*s = replacer
	}
}

func init() {
	err := envconfig.Process("", &Config)
	if err != nil {
		log.Errorf("Failed to load test env: %s", err.Error())
	}
	replaceIfEmpty(&Config.RepoDir, defaultRepoDir)
}
