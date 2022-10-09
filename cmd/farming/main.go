package farming

import (
	"github.com/MatroxMC/FS22ServerManager/internal/game"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

/*
 * This Farming struct is used to store the configuration of the game in config file
 *
 * Directory: Is the directory where the game is installed
 * Steam : true if the game is installed on steam, if steam is not started the process wait steam to start
 * Window : true if the game window is shown
 */

type Farming struct {
	Directory   string      `toml:"directory"`
	Steam       steam.Steam `toml:"steam"`
	Window      bool        `toml:"window"`
	RestartTime int         `toml:"restart_time"`
	game        game.Game
}

func (f Farming) Start(w *sync.WaitGroup) error {
	g, err := game.New(f.Directory, f.Steam, f.Window)
	if err != nil {
		return err
	}

	g.Signal = make(chan os.Signal, 1)

	g.NewHandler(func() error {
		go g.Restart(time.Duration(f.RestartTime)) //bad idea
		return nil
	}, game.HandleClosed)

	w.Add(1)
	go func() {
		defer w.Done()

		err := g.Start()
		if err != nil {
			panic(err)
		}
	}()

	f.game = *g

	return nil
}

func (f Farming) Stop() {
	err := f.game.Stop()
	if err != nil {
		log.Printf("Error while stopping the game: %s", err)
	}
}

func PrintInfo(g game.Game) {
	log.Printf(strings.Repeat("―", 50))
	log.Printf("Version : %s", g.Info.String)
	log.Printf("Steam : %t", g.Steam)
	log.Printf("Show Window : %t", g.ShowWindow)
	log.Printf("Directory : %s", g.Directory)
	log.Printf(strings.Repeat("―", 50))
}
