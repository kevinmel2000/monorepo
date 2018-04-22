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
	if !strings.Contains(dir, "/lab46/monorepo") {
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
		"../../../svc/",
		"../../../svc/bookapp",
		"../../../svc/bookapp/book",
		"../../../svc/bookapp/book/book.go",
		"../../../svc/bookapp/book/old",
		"../../../svc/bookapp/book/something.go",
		"../../../svc/rentapp/rent",
		"../../../svc/rentapp/rent/old",
		"../../../svc/rentapp/rent/something.go",
		"../../../experiment/grpc",
		"../../../experiment/grpc/svc1",
		"../../../experiment/grpc/svc2",
	}
	expect := []string{
		"../../../svc/bookapp",
		"../../../svc/rentapp",
	}

	dirs, err := FilterDir(files)
	if err != nil {
		t.Error(err)
		return
	}

	if len(dirs) != len(expect) {
		t.Error("dirs and expected length mismatch")
		return
	}

	for _, e := range expect {
		if _, ok := dirs[e]; !ok {
			t.Errorf("Dir %s is not exists", e)
			t.Error("Dirs:", dirs)
			return
		}
	}
}
