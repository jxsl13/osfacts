package common

import (
	"os/exec"
	"strings"
)

type Command struct {
	*exec.Cmd
	executable string
	args       []string
}

// OutputLines runs the command and returns the output (stdout) as lines
func (c *Command) OutputLines() ([]string, error) {
	b, err := c.Output()
	if err != nil {
		return nil, err
	}

	return SplitLines(strings.TrimSpace(string(b))), nil
}

func (c *Command) Output() ([]byte, error) {
	return c.Cmd.Output()

}

// OutputLines runs the command and returns the output (stdout) as lines
func (c *Command) Reset() {
	c.Cmd = exec.Command(c.executable, c.args...)
}

func GetCommand(executable string, args ...string) (cmd *Command, err error) {
	executablePath, err := exec.LookPath(executable)
	if err != nil {
		return nil, err
	}
	cmd = &Command{
		executable: executablePath,
		args:       args,
	}
	cmd.Reset()

	return cmd, nil
}
