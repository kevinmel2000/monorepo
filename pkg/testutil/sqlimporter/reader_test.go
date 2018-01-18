package sqlimporter

import "testing"

func TestGetFileList(t *testing.T) {
	fileList := []string{"test1.sql", "test2.sql"}
	list, err := getFileList("files")
	if err != nil {
		t.Error(err)
	}
	if len(list) != len(fileList) {
		t.Error("List of files is different")
	}
}
