package exec

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/lab46/example/gopkg/print"
)

var execPrint *print.Printer

func init() {
	execPrint = print.WithPrefix("[EXEC]")
}

// Cmd struct
type Cmd struct {
	name        string
	cmd         *exec.Cmd
	executedAt  time.Time
	mustSuccess bool
}

// Command for os/exec wrapper
func Command(name string, arg ...string) *Cmd {
	c := exec.Command(name, arg...)
	cmd := Cmd{
		cmd: c,
	}
	return &cmd
}

// MustSuccess state that the cmd is must success
func (c *Cmd) MustSuccess() {
	c.mustSuccess = true
}

// SetOutputToOS will set output to os output
func (c *Cmd) SetOutputToOS() {
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
}

// IgnoreOutput set output to nil
func (c *Cmd) IgnoreOutput() {
	c.cmd.Stdout = nil
	c.cmd.Stderr = nil
}

func (c *Cmd) isMustSuccess(err error) {
	if c.mustSuccess && err != nil {
		execPrint.Error(err.Error())
		os.Exit(1)
	}
}

// Output of os exec
func (c *Cmd) Output() ([]byte, error) {
	execPrint.Debug("Executing: [", c.name, c.cmd.Args, "]")
	t := time.Now()
	output, err := c.cmd.Output()
	elapsedTime := time.Since(t)
	execPrint.Debug("Command: [", c.name, c.cmd.Args, "] ~fin in:", fmt.Sprintf("%.3fs", elapsedTime.Seconds()))
	c.isMustSuccess(err)
	return output, err
}

// Run of os exec
func (c *Cmd) Run() error {
	execPrint.Debug("Executing: [", c.name, c.cmd.Args, "]")
	t := time.Now()
	err := c.cmd.Run()
	elapsedTime := time.Since(t)
	execPrint.Debug("Command: [", c.name, c.cmd.Args, "] ~fin in", fmt.Sprintf("%.3fs", elapsedTime.Seconds()))
	c.isMustSuccess(err)
	return err
}

// CombinedOutput for os exec
func (c *Cmd) CombinedOutput() ([]byte, error) {
	execPrint.Debug("Executing: [", c.name, c.cmd.Args, "]")
	t := time.Now()
	output, err := c.cmd.CombinedOutput()
	elapsedTime := time.Since(t)
	execPrint.Debug("Command: [", c.name, c.cmd.Args, "] ~fin in", fmt.Sprintf("%.3fs", elapsedTime.Seconds()))
	c.isMustSuccess(err)
	return output, err
}
