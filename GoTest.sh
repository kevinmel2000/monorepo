#!/bin/bash

filespath=()
# get files changed in the last commit
files="$(git diff-tree --no-commit-id --name-only -r ${DRONE_COMMIT_SHA})"
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
        cmd="go test -v $(go list ./... | grep -v /vendor/ | grep ${path})"
        $cmd
    fi
done