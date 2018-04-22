# Git-test 

Git-test is a CLI tools used to detect a diff in a service and run the task inside it.

## Service directory

Service directory need to be identified before it can detect changes. 

`git-test` using environment variables to control the configuration:

1. `GT_REPO_NAME`: Monorepo name
2. `GT_SERVICE_FOLDER`: Service directory of monorepo
3. `GT_REPO_DIR`: Directory of logistic monorepo
4. `GT_ENV`: Environment of git-test

All used environment variables can be found in [projectenv](https://github.com/lab46/monorepo/blob/master/tools/git-test/projectenv/projectenv.go) package

## Use Git at its core

`git-test` is using 2 git commands to detect changes

1. `git status -s` to detect change before commit is created

2. `git diff-tree -no-commit-id --name-only -r ${SHA1}` to detect changes in a given commit

All of git command above can be found in [git](https://github.com/lab46/monorepo/blob/master/tools/git-test/git/git.go) package 

## Running a task when changes is detected

When changes detected in a service, a task need to be defined via `task.yaml`. If the task is not detected, no task will run when changes detected.

### Task Runner

Task runner is defined as `yaml`. For example:

```yaml
test:
    - name: doing ls
      command: ls
      env: ['test']
    - name: doing ls -la
      command: ls -la
```

Test task is an array task, which each task have:
1. Name (`mandatory`)
    - Name is name of the task, this is will be printed by the `runner` so user know what task is running
2. Command (`mandatory`)
    - Command is the executeable command that `runner` can `exec`. All command output will be printed into `stdout`
3. Env (`optional`)
    - Env is environment of `git-test`, this is useful to run a command in a spesific environment like `test`.

## Command

Example of `git-test help`

```
git-test command line tools

Usage:
  git-test [command]

Available Commands:
  commit      commit will detect changes in a commit and test it
  diff        diff will detect diff before commit and test the changes
  help        Help about any command
  info        info command for git-test
  service     service will test a service in service directory

Flags:
  -h, --help      help for git-test
  -v, --verbose   sqlimporter verbose output

Use "git-test [command] --help" for more information about a command.
```

## Git test in action

Coming soon, will be a gif