package main

import (
	"fmt"
	"time"

	"github.com/lab46/example/pkg/print"
	"github.com/lab46/example/pkg/testutil/sqlimporter"
	"github.com/spf13/cobra"
)

var (
	// global variable from global flags
	VerboseFlag bool
	timeStart   time.Time
)

func initCMD() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sqlimporter [command]",
		Short: "sqlimporter command line tools",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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

func main() {
	rootCmd := initCMD()
	registerImporterCommand(rootCmd)
	rootCmd.Execute()
}

func registerImporterCommand(root *cobra.Command) {
	cmds := []*cobra.Command{
		{
			Use:   "import [driver] [dbname] [dsn] [directory]",
			Short: "import postgresql/mysql schema from directory",
			Args:  cobra.MinimumNArgs(3),
			Run: func(c *cobra.Command, args []string) {
				driver := args[0]
				dbName := args[1]
				dsn := args[2]
				dir := args[3]

				db, _, err := sqlimporter.CreateDB(driver, dbName, dsn)
				print.Fatal(err)
				err = sqlimporter.ImportSchemaFromFiles(db, dir)
				print.Fatal(err)
				print.Info("Successfully import schema from", dir)
			},
		},
		{
			Use:   "test [args]",
			Short: "test command for sqlimporter",
			Run: func(c *cobra.Command, args []string) {
				fmt.Println(args)
			},
		},
	}
	root.AddCommand(cmds...)
}
