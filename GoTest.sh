#!/bin/bash

# echo _/Users/Valge/Go/src/github.com/lab46/example/pkg/webserver | sed 's/\_\/Users\/Valge\/Go\/src\///'
# escape $GOPATH/src/ to \$GOPATH\/src\/
gopath=$(echo "$GOPATH/src/" | sed 's_/_\\/_g')
commitsha=$1
type=$2

# check commitsha from param
if [ "$commitsha" == "" ]; then
    echo "commit_sha is empty"
    exit 1
fi

# get all go packages in repo
# sometimes go list will print _/$GOPATH/src/project package instead $GOPATH/src/project/package
# need to trim _/$GOPATH/src before go test and go build can run
go_packages=()
files_to_test="$(go list ./... | grep -v /vendor/)"
for file in $files_to_test
do  
    go_package=${file}
    if [ "${file:0:1}" = "_" ]; then
        go_package=${file:1}
    fi

    # replacing $PWD/package to github.com/project/package
    # ex: Go/src/github.com/lab46/example/bookapp to github.com/lab46/example/bookapp
    go_package="$(echo $go_package | sed "s/$gopath//")"
    go_packages+=("${go_package}")
done

# for debug purposes
echo "RUNNING: ./GoTest.sh ${commitsha} ${type}"

files=""
# check if branch is the test, if yes then diff changes with master branch
if [ "$type" = "branch" ]; then
    files="$(git --no-pager diff --name-only $commitsha $(git merge-base $commitsha master))"
else
    # if diff, then test it with current git diff
    if [ "$commitsha" = "diff" ]; then
        files="$(git status -s | awk '{print $2}')"
    else   
        # get files changed in the last commit
        # will only get changed files in /*.go and not *.go
        files="$(git diff-tree --no-commit-id --name-only -r $commitsha | egrep "\/.+\.go")"    
    fi 
fi

if  [ "$files" = "" ]; then
    echo ">>> no files detected, exiting test..."
    exit 0
fi

# get every path in $files list in loop/iteration
# ex: "path1/file1.go path2/file1.go" into ["path1/file1.go", "path2/file1.go"]
filespath=()
for path in $files
do
    # split filepath and only get 1st string in array
    # ex: filepath/something.go, then only "filepath" is appended to filespath
    # TODO: Using sed instead of this, and remove *.go
    IFS='/' read -ra filepath <<< "$path"
    filespath+=("${filepath[0]}")
done

# get unique path from array of string
unique_path=($(echo "${filespath[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))

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

# check go pakcage available 
if [ "$(echo ${go_test_pkg[@]})" = "" ]; then
    echo ">>> no Go package detected, exiting test..."
fi

# do test and build
for test_pkg in ${go_test_pkg[@]}
do
    test="go test -v -race $test_pkg"
    $test
    build="go build -v $test_pkg/..."
    $build
done