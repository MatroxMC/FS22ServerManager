package game

import (
	"fmt"
	"github.com/MatroxMC/FS22ServerManager/internal/game/version"
	"github.com/MatroxMC/FS22ServerManager/internal/process"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"log"
	"os"
	"os/signal"
	"path"
)

type Version string
type Binary string

// This var save all game version
var versions = []version.Version{
	version.FS22{},
	version.FS19{},
}

type Game struct {
	Binary      Binary
	Steam       steam.Steam
	Version     version.Version
	Process     *process.Process
	ShowWindow  bool
	Directory   string
	HandleClose func(Game, error) error
	HandleStart func(Game) error
	Signal      chan os.Signal
	Killed      bool
}

var CurrentGame Game

func (g Game) Start() error {
	//Create a new process
	p, err := process.New(g.Binary.String())
	if err != nil {
		return err
	}

	g.Process = p   //set the process to the game
	CurrentGame = g //Set the current game to the game

	//If steam use steam
	if g.Steam {
		if !g.Steam.IsInstalled() {
			return fmt.Errorf("steam is not installed")
		}
	}

	if g.HandleStart != nil {
		err = g.HandleStart(g)
		if err != nil {
			return err
		}
	}

	// If steam not running the wait function is started and wait steam to run
	if !g.Steam.IsRunning() {
		log.Print("Steam is not running, waiting for steam to run")
		err := g.Steam.Wait()
		if err != nil {
			return err
		}
	}

	//Run the process with process package
	err = p.Run()
	if err != nil {
		return err
	}

	err = p.Wait()

	if g.HandleClose != nil {
		err := g.HandleClose(g, nil)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
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

func (g Game) Restart() error {
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

func New(directory string, version Version, steam steam.Steam, window bool) (*Game, error) {
	//check if version exist
	v, err := version.Find()
	if err != nil {
		return nil, err
	}

	//check if game binary exist
	binary := Binary(path.Join(directory, v.BinaryName()))
	err = binary.Exist()
	if err != nil {
		return nil, err
	}

	return &Game{
		Version:    v,
		Binary:     binary,
		Steam:      steam,
		Directory:  directory,
		ShowWindow: window,
	}, nil
}

func (v Version) Find() (version.Version, error) {
	for _, vv := range versions {
		for _, name := range vv.Names() {
			if name == string(v) {
				return vv, nil
			}
		}
	}

	return version.FS22{}, fmt.Errorf("version %s not found", v)
}

func (v Version) String() string {
	return string(v)
}

func (d Binary) Exist() error {
	return file.Exist(string(d))
}

func (d Binary) String() string {
	return string(d)
}
