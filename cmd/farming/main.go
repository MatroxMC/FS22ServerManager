package farming

import (
	"github.com/MatroxMC/FS22ServerManager/internal/game"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"log"
	"strings"
)

const (
	GameDir  = "D:\\Jeux\\Farming Simulator 22"
	ConfName = "config.toml"
)

type Farming struct {
	Directory game.Binary `toml:"directory"`
	Steam     steam.Steam `toml:"steam"`
	Window    bool        `toml:"window"`
}

func (f Farming) Start() (game.Game, error) {
	g, err := game.New(f.Directory.String(), f.Steam, f.Window)
	if err != nil {
		return game.Game{}, err
	}

	log.Printf(strings.Repeat("―", 50))
	log.Printf("Version : %s", g.Info.String)
	log.Printf("Steam : %t", g.Steam)
	log.Printf("Show Window : %t", g.ShowWindow)
	log.Printf("Binary : %s", g.Binary)
	log.Printf("Directory : %s", g.Directory)
	log.Printf(strings.Repeat("―", 50))

	g.HandleStart = func(game game.Game) error {
		game.Process.ShowWindow(game.ShowWindow) //Set the process parameter before start
		return nil
	}

	g.HandleClose = func(game game.Game, err error) error {
		if !game.Process.Killed {

			log.Println("Game closed, restarting in few seconds...")
			_, err := game.Restart()
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

	g, err = g.Start() //Run game instance
	if err != nil {
		return game.Game{}, err
	}

	err = g.Process.Wait()
	if err != nil {
		return game.Game{}, err
	}

	return *g, nil
}
