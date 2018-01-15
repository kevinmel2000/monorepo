#!/bin/bash

# escape $GOPATH/src/ to \$GOPATH\/src\/
gopath=$(echo "$GOPATH/src/" | sed 's_/_\\/_g')
commitsha=$1
type=$2

# check commitsha from param
if [ "$commitsha" == "" ]; then
    echo "commit_sha is empty"
    exit 1
fi

echo "RUNNING: ./GoTest.sh ${commitsha} ${type}"
filespath=()

files=""
# check if branch is the test, if yes then diff changes with master branch
if [ "$type" = "branch" ]; then
    files="$(git --no-pager diff --name-only $commitsha $(git merge-base $commitsha master))"
else
    # if diff, then test it with current git diff
    if [ "$commitsha" = "diff" ]; then
        files="$(git status -s | awk '{print $2}' | egrep "\/.+\.go")"
        # files="$(git diff --name-only | egrep "\/.+\.go")"
    else   
        # get files changed in the last commit
        # will only get changed files in /*.go and not *.go
        files="$(git diff-tree --no-commit-id --name-only -r $commitsha | egrep "\/.+\.go")"    
    fi 
fi

if  [ "$files" = "" ]; then
    echo "no Go files detected, exiting test..."
    exit 0
fi

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
        # sometimes go list will print _/$GOPATH/src/project package instead $GOPATH/src/project/package
        # need to trim _/$GOPATH/src before go test and go build can run
        files_to_test="$(go list ./... | grep -v /vendor/ | grep ${path})"
        for file in $files_to_test
        do
            testfile=${file}
            if [ "${file:0:1}" = "_" ]; then
                testfile=${file:1}
            fi

            # replacing $PWD/package to github.com/project/package
            # ex: Go/src/github.com/lab46/example/bookapp to github.com/lab46/example/bookapp
            testfile="$(echo $testfile | sed "s/$gopath//")"
            # run test and build
            test="go test -v $testfile"
            $test
            build="go build -v $testfile
            $build"
        done
    fi
done