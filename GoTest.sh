#!/bin/bash
# set -e

commitsha=$1
echo "DEBUG: ./GetTest.sh ${commitsha}"
filespath=()
# get files changed in the last commit
# will only get changed files in /*.go and not *.go
files="$(git diff-tree --no-commit-id --name-only -r $commitsha | egrep "\/.+\.go")"
# get every path in $files list in loop/iteration
# ex: "path1/file1.go path2/file1.go" into ["path1/file1.go", "path2/file1.go"]
for path in $files
do
    # split filepath and only get 1st string in array
    # ex: filepath/something.go, then only "filepath" is appended to filespath
    IFS='/' read -ra filepath <<< "$path"
    filespath+=("${filepath[0]}")
done

# get unique path from array of string
unique_path=($(echo "${filespath[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))
for path in ${unique_path[@]}
do
    # only test if path is directory
    if [ -d "$path" ]; then
        test="go test -v $(go list ./... | grep -v /vendor/ | grep ${path})"
        $test
        build="go build -v $(go list ./... | grep -v /vendor/ | grep ${path})"
        $build
    fi
done