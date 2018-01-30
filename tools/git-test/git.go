package main

import (
	"strings"

	"github.com/lab46/example/gopkg/exec"
	"github.com/lab46/example/gopkg/print"
)

// get diff list exec 'git status -s'
// but all deleted files will be skipped by this function
func gitDiffList() ([]string, error) {
	// check changes using git status command
	cmd := exec.Command("git", "status", "-s")
	cmd.MustSuccess()
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	filesStatus := strings.Split(string(output), "\n")

	var changedFiles []string
	for _, file := range filesStatus {
		// continue if empty string
		if len(file) == 0 {
			continue
		}

		// sometimes we got "?? files" or " M" so 1st string need to be checked
		awk := 1
		if file[:1] == " " {
			awk++
		}
		f := strings.Split(file, " ")
		if len(f) < 2 {
			continue
		}
		status := f[awk-1]
		// skip deleted files
		if status == "D" {
			continue
		}
		changedFiles = append(changedFiles, f[awk])
	}
	return changedFiles, nil
}

func gitSHAList(SHA1 string) ([]string, error) {
	// git diff-tree --no-commit-id --name-only -r e38170994b15aa70016e3f57e94a2110c36842a2
	cmd := exec.Command("git", "diff-tree", "-no-commit-id", "--name-only", "-r", SHA1)
	cmd.MustSuccess()
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	print.Debug(string(output))

	files := strings.Split(string(output), "\n")
	// if only commit SHA available
	if len(files) == 1 {
		return nil, nil
	}

	var changedFiles []string
	for _, file := range files[1:] {
		changedFiles = append(changedFiles, file)
	}
	return changedFiles, err
}
