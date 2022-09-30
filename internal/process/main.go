package process

import (
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"os"
	"os/exec"
	"syscall"
)

const (
	RUNNING = 0
	STOPPED = 1
)

type Process struct {
	Executable string
	Cmd        exec.Cmd
	Signal     chan os.Signal
	Killed     bool
}

type Status int

func signal() {

}

func New(path string) (*Process, error) {
	err := file.Exist(path)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path)

	return &Process{
		Executable: path,
		Cmd:        *cmd,
	}, nil
}

func (p *Process) ShowWindow(d bool) {
	p.Cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: !d,
	}
}

func (p *Process) Start() error {
	return p.Cmd.Start()
}

func (p *Process) Stop() error {

	return nil
}

func (p *Process) Running() bool {
	return p.Cmd.ProcessState != nil
}

func (p *Process) Wait() error {
	return p.Cmd.Wait()
}

func (p *Process) Kill() error {
	if cmd := p.Cmd.Process; cmd != nil {
		p.Killed = true

		err := cmd.Kill()
		if err != nil {
			return err
		}
	}
	return nil
}
