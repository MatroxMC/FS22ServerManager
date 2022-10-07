package game

import (
	"fmt"
	"github.com/MatroxMC/FS22ServerManager/internal/process"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"log"
	"os"
	"os/signal"
	"path"
)

type Binary string

type Game struct {
	Binary      Binary
	Steam       steam.Steam
	Info        Info
	Process     *process.Process
	ShowWindow  bool
	Directory   string
	HandleClose func(Game, error) error
	HandleStart func(Game) error
	Signal      chan os.Signal
	Killed      bool
}

var CurrentGame Game

func (g Game) Start() (*Game, error) {
	//Create a new process
	p, err := process.New(g.Binary.String())
	if err != nil {
		return &Game{}, err
	}

	g.Process = p   //set the process to the game
	CurrentGame = g //Set the current game to the game

	//If steam use steam
	if g.Steam {
		if !g.Steam.IsInstalled() {
			return &Game{}, fmt.Errorf("steam is not installed")
		}
	}

	if g.HandleStart != nil {
		err = g.HandleStart(g)
		if err != nil {
			return &Game{}, err
		}
	}

	// If steam not running the wait function is started and wait steam to run
	if !g.Steam.IsRunning() {
		log.Print("Steam is not running, waiting for steam to run")
		err := g.Steam.Wait()
		if err != nil {
			return &Game{}, err
		}
	}

	//Run the process with process package
	err = p.Run()
	if err != nil {
		return &Game{}, err
	}

	return &g, nil
}

func (g Game) Init() (*Game, error) {
	g.Signal = make(chan os.Signal, 1) //Create a new channel for signal

	g.handleSignal(func() {
		//I use Current game because this is in go routine (TODO: Clean this in the future)
		if CurrentGame.Process != nil {
			_ = CurrentGame.Process.Stop()
		}
	})

	return &g, nil
}

func (g Game) Restart() (*Game, error) {
	//Check if the process running and kill it
	if g.Process.Running() {
		_ = g.Process.Stop()
	}

	//Run the process again
	return g.Start()
}

// This method handle the signal from main process and kill the game process
func (g Game) handleSignal(f ...func()) {
	signal.Notify(g.Signal)

	go func() {
		_ = <-g.Signal
		if len(f) > 0 {
			//Handle all function
			for _, fn := range f {
				fn()
			}
		}

		signal.Stop(g.Signal)
		os.Exit(0) //Exit the process with code 0 (no error)
	}()
}

func New(directory string, steam steam.Steam, window bool) (*Game, error) {
	//check if game binary exist
	binary := Binary(path.Join(directory))
	err := binary.Exist()
	if err != nil {
		return nil, err
	}

	return &Game{
		Info:       DefaultInfo(),
		Binary:     binary,
		Steam:      steam,
		Directory:  directory,
		ShowWindow: window,
	}, nil
}

func (d Binary) Exist() error {
	return file.Exist(string(d))
}

func (d Binary) String() string {
	return string(d)
}
