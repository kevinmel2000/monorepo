package task

import (
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
