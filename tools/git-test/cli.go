package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/lab46/monorepo/gopkg/exec"
	"github.com/lab46/monorepo/gopkg/print"
	"github.com/lab46/monorepo/tools/git-test/git"
	"github.com/lab46/monorepo/tools/git-test/projectenv"
	"github.com/lab46/monorepo/tools/git-test/repo"
	"github.com/lab46/monorepo/tools/git-test/runner"
)

var (
	// VerboseFlag for debugging
	VerboseFlag bool
	timeStart   time.Time
)

func checkDependencies() {
	// check if git is available
	cmd := exec.Command("git", "version")
	cmd.MustSuccess()
	cmd.IgnoreOutput()
	cmd.Run()

	// check if go is available
	cmd = exec.Command("go", "version")
	cmd.MustSuccess()
	cmd.IgnoreOutput()
	cmd.Run()

	// check if dep is available
	cmd = exec.Command("dep", "version")
	cmd.MustSuccess()
	cmd.IgnoreOutput()
	cmd.Run()
}

func initCMD() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "git-test [command]",
		Short: "git-test command line tools",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			checkDependencies()
			print.SetDebug(VerboseFlag)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			elapsedTime := time.Since(timeStart).Seconds()
			print.Info("Command <<", cmd.CommandPath(), args, ">> running in", fmt.Sprintf("%.3fs", elapsedTime))
		},
	}
	// add flags
	rootCmd.PersistentFlags().BoolVarP(&VerboseFlag, "verbose", "v", false, "sqlimporter verbose output")
	timeStart = time.Now()
	return rootCmd
}

func registerGitTestCommand(root *cobra.Command) {
	cmds := []*cobra.Command{
		{
			Use:   "diff",
			Short: "diff will detect diff before commit and test the changes",
			Args:  cobra.MaximumNArgs(1),
			Run: func(c *cobra.Command, args []string) {
				print.Error(changeToRepodir())
				list, err := git.DiffList()
				print.Error(err)
				print.Error(processDirs(list))
			},
		},
		{
			Use:   "commit [commitsha1]",
			Short: "commit will detect changes in a commit and test it",
			Args:  cobra.MinimumNArgs(1),
			Run: func(c *cobra.Command, args []string) {
				print.Error(changeToRepodir())
				sha1 := args[0]
				list, err := git.SHAList(sha1)
				print.Error(err)
				print.Fatal(processDirs(list))
			},
		},
		{
			Use:   "service [servicename ...]",
			Short: "service will test a service in service directory",
			Args:  cobra.MinimumNArgs(1),
			Run: func(c *cobra.Command, args []string) {
				serviceFolder := repo.ServiceFolder()
				services := make([]string, len(args))
				for key, val := range args {
					services[key] = path.Join(serviceFolder, val)
				}
				print.Fatal(processDirs(services))
			},
		},
		{
			Use:   "info [args]",
			Short: "info command for git-test",
			Run: func(c *cobra.Command, args []string) {
				repoName := projectenv.Config.RepoName
				serviceDir := projectenv.Config.ServiceFolder
				repoDir := projectenv.Config.RepoDir

				print.Info("REPO_NAME:", repoName)
				print.Info("REPO_DIR:", repoDir)
				print.Info("SERVICE_FOLDER:", serviceDir)
			},
		},
	}
	root.AddCommand(cmds...)
}

// changeToRepodir for changedir to root repo path, and make git-test work on relative repo path
func changeToRepodir() error {
	// change dir into root path of repository
	repoDir, err := repo.GetRepoDir()
	if err != nil {
		return err
	}
	return os.Chdir(repoDir)
}

func changeToServiceDir() error {
	// change dir into service path of repository
	serviceDir, err := repo.GetServiceDir()
	if err != nil {
		return err
	}
	return os.Chdir(serviceDir)
}

func processDirs(list []string) error {
	dirs, err := repo.FilterDir(list)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		switch dir.Type {
		case repo.ServiceType:
			err := runner.TriggerServiceRunner(dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
