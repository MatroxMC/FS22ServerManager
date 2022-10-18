package server

import (
	"fmt"
	"github.com/MatroxMC/FS22ServerManager/internal/event"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"os"
	"os/exec"
	"syscall"
)

const (
	ProgramInit    event.HandleType = 0x0
	ProgramRun     event.HandleType = 0x1
	ProgramStarted event.HandleType = 0x2
	ProgramStopped event.HandleType = 0x3
	ProgramClosed  event.HandleType = 0x4
)

type Program struct {
	Executable string
	Directory  string
	ShowWindow bool
	Cmd        exec.Cmd
	Steam      steam.Steam
	Handler    event.Handler
	onShutdown bool
}

func (p *Program) Init() error {
	e := exec.Command(p.Executable)
	e.Dir = p.Directory
	e.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: !p.ShowWindow,
	}

	p.Cmd = *e

	err := p.Handler.HandleFunc(ProgramInit)
	if err != nil {
		return err
	}

	return nil
}

func (p *Program) Start() error {
	err := p.Init()
	if err != nil {
		return err
	}

	//Start the program
	err = p.Cmd.Start()
	if err != nil {
		return err
	}

	err = p.Handler.HandleFunc(ProgramRun)
	if err != nil {
		return err
	}

	go func() {
		//Wait for the program to exit
		err = p.Cmd.Wait()
		if err != nil {
			err = p.Handler.HandleFunc(ProgramStopped)
			if err != nil {
				return
			}
		}

		err = p.Handler.HandleFunc(ProgramClosed)
		if err != nil {
			return
		}
	}()

	return nil
}

func (p *Program) Stop() error {
	if p.onShutdown {
		return nil
	}

	p.onShutdown = true

	if p.Cmd.Process == nil {
		return nil
	}

	err := p.Cmd.Process.Kill()
	if err != nil {
		return err
	}

	return nil
}

func NewProgram(d string, e string, s steam.Steam, window bool) (Program, error) {

	ee := fmt.Sprintf("%s\\%s", d, e)

	if _, err := os.Stat(ee); os.IsNotExist(err) {
		return Program{}, fmt.Errorf("%s does not exist", e)
	}

	return Program{
		Executable: ee,
		Directory:  d,
		ShowWindow: window,
		Steam:      s,
	}, nil
}
