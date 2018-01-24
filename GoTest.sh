#!/bin/bash

# escape $GOPATH/src/ to \$GOPATH\/src\/
gopath=$(echo "$GOPATH/src/" | sed 's_/_\\/_g')
commitsha=$1
type=$2

# check commitsha from param
if [[ "$commitsha" == "" ]]; then
    echo "commit_sha is empty"
    exit 1
fi

# for debug purposes
echo "RUNNING: ./GoTest.sh ${commitsha} ${type}"

files=""
# check if branch is the test, if yes then diff changes with master branch
if [[ "$type" = "branch" ]]; then
    files="$(git --no-pager diff --name-only $commitsha $(git merge-base $commitsha master))"
else    
    # if diff, then test it with current git diff
    if [[ "$commitsha" = "diff" ]]; then
        files="$(git status -s | awk '{print $2}')"
    else   
        # get files changed in the last commit
        # will only get changed files in /*.go and not *.go
        files="$(git diff-tree --no-commit-id --name-only -r $commitsha | egrep "\/.+\.go")"    
    fi 
fi

if  [[ "$files" = "" ]]; then
    echo ">>> no files detected, exiting test..."
    exit 0
fi

filespath=()
# get every path in $files list in loop/iteration
# ex: "path1/file1.go path2/file1.go" into ["path1/file1.go", "path2/file1.go"]
for path in $files
do
    # trim filepath, only get the package name
    filespath+=("$(echo $path | rev | cut -d"/" -f2-  | rev)")
done
# get unique path from array of string
unique_path=($(echo "${filespath[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))

# get all go packages in repo
# sometimes go list will print _/$GOPATH/src/project package instead $GOPATH/src/project/package
# need to trim _/$GOPATH/src before go test and go build can run
# example: echo _/Users/Valge/Go/src/github.com/lab46/example/pkg/webserver | sed 's/\_\/Users\/Valge\/Go\/src\///'
go_packages="$(go list ./... | grep -v /vendor/ | sed "s/\_$gopath//")"

# looks for go test path
# need to improve this, very2 slow
go_test_pkg=()
for path in ${unique_path[@]}
do
    for package in ${go_packages[@]}
    do
        pack="$(echo $package | grep $path)"
        go_test_pkg+=("$pack")
    done
done

# check if a list of go pakcage is available 
if [[ "$(echo ${go_test_pkg[@]})" = "" ]]; then
    echo ">>> no Go package detected, exiting test..."
    exit 0
fi

#set exit when test failed
set -e
# do test and build
for test_pkg in ${go_test_pkg[@]}
do
    test="go test -v -race -parallel 2 $test_pkg"
    $test
    build="go build -v $test_pkg/..."
    $build
done