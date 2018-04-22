package runner

import (
	"os"
	"testing"

	"github.com/lab46/monorepo/tools/git-test/repo"
)

func TestTriggerServiceRunner(t *testing.T) {
	dir := repo.Dir{
		Name: "files/testfile/task",
		Type: repo.ServiceType,
	}

	repoDir, err := repo.GetRepoDir()
	if err != nil {
		t.Error(err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Error(err)
	}

	err = TriggerServiceRunner(dir)
	if err != nil {
		t.Error(err)
	}
}
