package runner

import (
	"testing"

	"github.com/lab46/example/tools/git-test/repo"
)

func TestTriggerServiceRunner(t *testing.T) {
	dir := repo.Dir{
		Name: "../../../files/testfile/task",
		Type: repo.ServiceType,
	}
	err := TriggerServiceRunner(dir)
	if err != nil {
		t.Error(err)
	}
}
