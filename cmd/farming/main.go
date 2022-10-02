package farming

import (
	"github.com/MatroxMC/FS22ServerManager/internal/game"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"log"
)

const (
	GameDir  = "D:\\Jeux\\Farming Simulator 22"
	ConfName = "config.toml"
)

type Farming struct {
	Directory game.Binary  `toml:"directory"`
	Steam     steam.Steam  `toml:"steam"`
	Version   game.Version `toml:"version"`
	Window    bool         `toml:"window"`
}

func (f Farming) Start() (game.Game, error) {
	g, err := game.New(f.Directory.String(), f.Version, f.Steam, f.Window)
	if err != nil {
		return game.Game{}, err
	}

	log.Printf("-------- Server Manager (FS22-FS19) ---------")
	log.Printf("Version : %s", g.Version.String())
	log.Printf("Steam : %t", g.Steam)
	log.Printf("Show Window : %t", g.ShowWindow)
	log.Printf("Binary : %s", g.Binary)
	log.Printf("Directory : %s", g.Directory)
	log.Printf("---------------------------------------------")

	g.HandleStart = func(game game.Game) error {
		game.Process.ShowWindow(game.ShowWindow) //Set the process parameter before start
		return nil
	}

	g.HandleClose = func(game game.Game, err error) error {
		if !game.Process.Killed {

			log.Println("Game closed, restarting in few seconds...")
			err := game.Restart()
			if err != nil {
				return err
			}
		}
		return nil
	}

	g, err = g.Init()
	if err != nil {
		return game.Game{}, err
	}

	err = g.Start() //Run game instance
	if err != nil {
		return game.Game{}, err
	}

	return *g, nil
}
