package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/lab46/monorepo/tools/git-test/projectenv"

	"github.com/lab46/monorepo/gopkg/env"
	"github.com/lab46/monorepo/gopkg/exec"
	"github.com/lab46/monorepo/gopkg/print"
	"gopkg.in/yaml.v2"
)

// TaskFile is a task file name
const TaskFile = "task.yaml"

// GitTestEnv is env var for git-test
const GitTestEnv = "GIT_TEST_ENV"

// Tasks struct
type Tasks struct {
	T []Task `yaml:"test"`
}

// Task struct
type Task struct {
	Name string   `yaml:"name"`
	Run  string   `yaml:"run"`
	Env  []string `yaml:"env"`
}

// ReadTasksFromFile return task defined in yml file
func ReadTasksFromFile(filename string) (Tasks, error) {
	t := Tasks{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return t, err
	}
	err = yaml.Unmarshal(content, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func convertToEnvMap(envs []string) map[string]string {
	e := make(map[string]string)
	for _, env := range envs {
		e[env] = env
	}
	return e
}

func getCurrentEnv() string {
	e := env.Getenv(GitTestEnv)
	if e == "" {
		// default is dev
		e = "dev"
	}
	return e
}

// DoTasks execute all tasks
func DoTasks(t Tasks) error {
	currentEnv := projectenv.Config.Env
	for _, task := range t.T {
		// check environment, skip if not included
		if len(task.Env) > 0 {
			envs := convertToEnvMap(task.Env)
			if _, ok := envs[currentEnv]; !ok {
				print.Debug(fmt.Sprintf("skipping %s. Env: %s", task.Run, currentEnv))
				continue
			}
		}
		// check if command is not available
		if task.Run == "" || task.Run == " " {
			continue
		}

		print.Info(fmt.Sprintf("[TASK] %s", task.Name))
		commands := strings.Split(task.Run, " ")
		cmd := exec.Command(commands[0], commands[1:]...)
		cmd.SetOutputToOS()
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// IsTaskFileExsits check if TaskFile is exists within the directory
func IsTaskFileExsits(dir string) (bool, error) {
	p := path.Join(dir, TaskFile)
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsTaskFileExistsInCurrentDir check if task file is exsits in current directotry
func IsTaskFileExistsInCurrentDir() (bool, error) {
	_, err := os.Stat(TaskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
