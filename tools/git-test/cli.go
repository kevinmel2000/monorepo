package main

import (
	"fmt"
	"time"

	"github.com/lab46/example/gopkg/exec"
	"github.com/lab46/example/gopkg/print"
	"github.com/lab46/example/tools/git-test/repo"
	"github.com/lab46/example/tools/git-test/runner"
	"github.com/spf13/cobra"
)

var (
	// global variable from global flags
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
			// checkDependencies()
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
				list, err := gitDiffList()
				print.Error(err)
				err = processDirs(list)
				print.Error(err)
			},
		},
		{
			Use:   "commit [commitsha1]",
			Short: "commit will detect changes in a commit and test it",
			Args:  cobra.MinimumNArgs(1),
			Run: func(c *cobra.Command, args []string) {
				sha1 := args[0]
				list, err := gitSHAList(sha1)
				print.Error(err)
				processDirs(list)
				print.Error(err)
			},
		},
		{
			Use:   "test [args]",
			Short: "test command for git-test",
			Run: func(c *cobra.Command, args []string) {
				print.Info("ARGS:", args)
			},
		},
	}
	root.AddCommand(cmds...)
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
