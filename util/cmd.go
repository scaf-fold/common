package util

import (
	"os/exec"
	"sync"
)

type Cmd struct {
	sync.Mutex
}

func NewCmd() *Cmd {
	return &Cmd{}
}

func (c *Cmd) ShellExec(cmd string) ([]byte, error) {
	if cmd != "" {
		exe := exec.Command("sh", "-c", cmd)
		return exe.Output()
	}
	return nil, nil
}

func (c *Cmd) SShellExec(cmd string) ([]byte, error) {
	if cmd != "" {
		c.Lock()
		defer c.Lock()
		exe := exec.Command("sh", "-c", cmd)
		return exe.Output()
	}
	return nil, nil
}
