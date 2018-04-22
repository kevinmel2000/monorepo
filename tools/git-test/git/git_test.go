package git_test

import (
	"testing"

	"github.com/lab46/monorepo/tools/git-test/git"
)

func TestDiffList(t *testing.T) {
	list, err := git.DiffList()
	if err != nil {
		t.Error(t)
	}
	if list == nil {
		t.Error("list is nil")
	}
}

func TestSHAList(t *testing.T) {
	// this is sha1 from logistic repo
	sha1 := "5e25211c9abe1c4bb2ea19add4b93ef132ae074c"
	list, err := git.SHAList(sha1)
	if err != nil {
		t.Error(err)
		return
	}
	if len(list) == 0 {
		t.Error("list length is 0")
	}
}
