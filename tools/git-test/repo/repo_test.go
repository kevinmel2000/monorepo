package repo

import (
	"strings"
	"testing"
)

func TestGetRepoDir(t *testing.T) {
	dir, err := GetRepoDir()
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(dir, "/lab46/example") {
		t.Errorf("Dir is not as expected: %s", dir)
	}
}

func TestGetServiceList(t *testing.T) {
	s, err := ServiceList()
	if err != nil {
		t.Error(err)
	}
	if len(s) == 0 {
		t.Error("expecting more than 1 service")
	}
}

func TestFilterDir(t *testing.T) {
	files := []string{
		"something",
		// this need to be fixed, files is not tracked
		"../../../.gitignore",
		"../../../files/",
		"../../../exservice/",
		"../../../exservice/bookapp",
		"../../../exservice/bookapp/book",
		"../../../exservice/bookapp/book/book.go",
		"../../../exservice/bookapp/book/old",
		"../../../exservice/bookapp/book/something.go",
		"../../../exservice/rentapp/rent",
		"../../../exservice/rentapp/rent/old",
		"../../../exservice/rentapp/rent/something.go",
	}
	expect := []string{
		"../../../exservice/bookapp",
		"../../../exservice/rentapp",
	}

	dirs, err := FilterDir(files)
	if err != nil {
		t.Error(err)
	}

	for _, e := range expect {
		if _, ok := dirs[e]; !ok {
			t.Errorf("Dir %s is not exists", e)
			t.Error("Dirs:", dirs)
			return
		}
	}
}
