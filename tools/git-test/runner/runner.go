package runner

import (
	"fmt"
	"os"

	"github.com/lab46/example/gopkg/print"

	"github.com/lab46/example/tools/git-test/repo"
	"github.com/lab46/example/tools/git-test/task"
)

// TriggerServiceRunner run service model task
func TriggerServiceRunner(dir repo.Dir) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	print.Info(fmt.Sprintf("[RUNNER] Switching to %s", dir.Name))
	err = os.Chdir(dir.Name)
	if err != nil {
		return err
	}

	// return if task file is not exists within the current path
	if exists, err := task.IsTaskFileExistsInCurrentDir(); err != nil {
		return err
	} else if !exists {
		return nil
	}

	ts, err := task.ReadTasksFromFile(task.TaskFile)
	if err != nil {
		return err
	}
	err = task.DoTasks(ts)
	if err != nil {
		return err
	}

	err = os.Chdir(currentDir)
	return err
}
