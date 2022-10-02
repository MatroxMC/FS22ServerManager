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
		var run bool = false //If the game is running
		var msg bool = false //If the message has been sent

		for !run {
			run = f.Steam.IsRunning()
			if !run {
				if !msg {
					log.Println("Try to detect Steam... (If you don't have steam, please set steam to false in the config file)")
					msg = true
				}

				time.Sleep(5 * time.Second)
			} else {
				pid, _ := f.Steam.GetPID()
				user, _ := f.Steam.GetAutoLoginUser()
				user = strings.Title(user)
				log.Printf("%s is logged in Steam with PID %o", user, pid)
			}
		}
	}

	err = g.Start()
	if err != nil {
		return game.Game{}, err
	}

	err = g.Wait()
	if err != nil {
		return game.Game{}, err
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
