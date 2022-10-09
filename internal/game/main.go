package game

import (
	"fmt"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	HandleClosed HandleType = 0x0
	HandleStart  HandleType = 0x1
	HandleStop   HandleType = 0x2
)

type HandleType int
type HandleFunction func() error

type Game struct {
	Steam       steam.Steam
	Info        Info
	Cmd         exec.Cmd
	ShowWindow  bool
	Directory   string
	Signal      chan os.Signal
	HandledFunc map[HandleType]HandleFunction
	Test        bool
}

type Info struct {
	Binary string
	Names  []string
	String string
}

func New(directory string, steam steam.Steam, window bool) (*Game, error) {

	//Check if the directory exist
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return &Game{}, fmt.Errorf("directory %s does not exist", directory)
	}

	//Return a clean Game struct
	return &Game{
		Info:       DefaultInfo(),
		Steam:      steam,
		Directory:  directory,
		ShowWindow: window,
	}, nil
}

func (g *Game) Start() error {
	if _, err := os.Stat(g.BinaryPath()); os.IsNotExist(err) {
		return fmt.Errorf("binary %s does not exist", g.BinaryPath())
	}

	e := exec.Command(g.BinaryPath())
	e.Dir = g.Directory
	e.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: !g.ShowWindow,
	}

	g.Cmd = *e

	err := g.handleFunction(HandleStart)
	if err != nil {
		return err
	}

	err = g.Cmd.Start()
	if err != nil {
		return err
	}

	err = g.Cmd.Wait()
	if err != nil {
		err := g.handleFunction(HandleClosed)
		if err != nil {
			return err
		}
	}

	err = g.handleFunction(HandleStop)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Restart(s time.Duration) error {
	log.Print("Restarting in ", s.Seconds(), " seconds")
	if err := g.Stop(); err != nil {
		return err
	}

	time.Sleep(s * time.Second)

	return g.Start()
}

func (g *Game) Stop() error {
	if g.Cmd.Process == nil {
		return nil
	}

	//check if app is running
	if err := g.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
		return nil
	}

	err := g.Cmd.Process.Kill()
	if err != nil {
		log.Println(err)
	}

	log.Printf("Killed %s", g.Info.String)

	return nil
}

func (g *Game) handleFunction(h HandleType) error {
	for t, f := range g.HandledFunc {
		if t == h {
			return f()
		}
	}

	return nil
}

func (g *Game) NewHandler(f HandleFunction, h HandleType) {

	if g.HandledFunc == nil {
		g.HandledFunc = make(map[HandleType]HandleFunction)
	}

	g.HandledFunc[h] = f
}

func (g *Game) BinaryPath() string {
	return fmt.Sprintf("%s\\%s", g.Directory, g.Info.Binary)
}

func DefaultInfo() Info {
	return Info{
		Binary: "dedicatedServer.exe",
		Names: []string{
			"Farming Simulator 22",
			"22",
			"FS22",
		},
		String: "Farming Simulator 22",
	}
}
