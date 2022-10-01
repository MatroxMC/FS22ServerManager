package farming

import (
	"github.com/MatroxMC/FS22ServerManager/internal/game"
	"github.com/MatroxMC/FS22ServerManager/internal/process"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"log"
	"strings"
	"time"
)

const (
	GameDir  = "D:\\Jeux\\Farming Simulator 22"
	ConfName = "config.toml"
)

type Farming struct {
	Directory  game.Binary  `toml:"directory"`
	Steam      steam.Steam  `toml:"steam"`
	Version    game.Version `toml:"version"`
	running    process.Status
	ShowWindow bool `toml:"show_window"`
}

func (f Farming) Start() (game.Game, error) {
	g, err := f.init()
	if err != nil {
		return game.Game{}, err
	}

	//If steam use steam
	if f.Steam {
		for err := f.Steam.IsRunning(); err != nil; {
			log.Println("Steam is not running, waiting...")
			time.Sleep(5 * time.Second)

			err = f.Steam.IsRunning() //Retry to check if steam is running or not
		}

		user, _ := f.Steam.ActiveUser()
		log.Println("Active user on steam:", strings.ToUpper(user))
	}

	return g, nil
}

func (f Farming) init() (game.Game, error) {
	g, err := game.New(f.Directory.String(), f.Version, f.Steam, f.ShowWindow)
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

	g.OnStart = func(game game.Game) error {
		game.Process.ShowWindow(game.ShowWindow)
		return nil
	}

	g.OnClose = func(game game.Game, err error) error {
		if !game.Process.Killed {
			log.Println("Game closed, restarting in few seconds...")
			err := game.Restart()
			if err != nil {
				return err
			}
		}

		return nil
	}

	return *g, nil
}
