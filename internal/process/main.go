package process

import (
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"os/exec"
	"syscall"
	"time"
)

type Process struct {
	Executable string
	Cmd        *exec.Cmd
}

func NewProcess(p string) (*Process, error) {
	err := file.Exist(p)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(p)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false, //hide console
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	return &Process{
		Executable: p,
		Cmd:        cmd,
	}, nil
}

func (p *Process) Start() error {
	err := p.Cmd.Start()

	return err
}

func (p *Process) Stop() error {
	//send CTRL + C
	err := p.Cmd.Process.Signal(syscall.SIGKILL)
	if err != nil {
		time.Sleep(5 * time.Second)
		err = p.Cmd.Process.Kill()
		if err != nil {
			return err
		}
	}

	return nil
}
