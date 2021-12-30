package common

import (
	"os/exec"
	"strings"
)

type Command struct {
	executable string
	args       []string
}

// OutputLines runs the command and returns the output (stdout) as lines
func (c *Command) OutputLines(args ...string) ([]string, error) {
	b, err := c.Output(args...)
	if err != nil {
		return nil, err
	}
	stripped := strings.TrimSpace(string(b))
	lines := SplitLines(stripped)
	return lines, nil
}

func (c *Command) Output(args ...string) ([]byte, error) {
	if len(args) > 0 {
		c.args = args
	}

	cmd := exec.Command(c.executable, c.args...)
	return cmd.Output()

}

func GetCommand(executable string, args ...string) (cmd *Command, err error) {
	executable, err = exec.LookPath(executable)
	if err != nil {
		return nil, err
	}
	cmd = &Command{
		executable: executable,
		args:       args,
	}

	return cmd, nil
}
