package task

import (
	"os"
	"testing"
)

func TestDoTask(t *testing.T) {
	testfile := "../../../files/testfile/task/task.yaml"
	ts, err := ReadTasksFromFile(testfile)
	if err != nil {
		t.Error(err)
	}
	err = DoTasks(ts)
	if err != nil {
		t.Error(err)
	}
}

func TestIsTaskFileExists(t *testing.T) {
	testfile := "../../../files/testfile/task"
	exist, err := IsTaskFileExsits(testfile)
	if err != nil {
		t.Error(err)
		return
	}
	if !exist {
		t.Error("file should be exists")
	}
}

func TestIsTaskFileExistsInCurrentDir(t *testing.T) {
	testdir := "../../../files/testfile/task"
	err := os.Chdir(testdir)
	if err != nil {
		t.Error(err)
		return
	}
	exist, err := IsTaskFileExistsInCurrentDir()
	if err != nil {
		t.Error(err)
	}
	if !exist {
		t.Error("file should be exist")
	}
}
