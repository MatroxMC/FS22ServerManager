package farming

import (
	"github.com/MatroxMC/FS22ServerManager/internal/game"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"github.com/kataras/golog"
	"os"
	"sync"
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

func (f *Farming) Init() error {
	g, err := game.Npoew(f.Directory, f.Steam, f.Window)
	if err != nil {
		return err
	}

	f.game = *g
	f.game.Signal = make(chan os.Signal, 1)
	return nil
}

func (f *Farming) Start(w *sync.WaitGroup) error {
	err := f.Init()
	if err != nil {
		return err
	}

	f.game.NewHandler(func() error {
		golog.Info("Server Manager started")
		return nil
	}, game.HandleStart)

	f.game.NewHandler(func() error {
		if !f.game.Killed {
			go f.game.Restart() //restart the game after 5 seconds
		}
		return nil
	}, game.HandleClosed)

	if f.Steam {
		if !f.Steam.IsRunning() {
			golog.Info("Steam is not running, waiting for steam to start")
			s := f.Steam.WaitForRunning()
			if s == steam.StatusExited {
				return nil
			}

			golog.Info("Steam is running")
		}
	}

	w.Add(1)
	go func() {
		defer w.Done()

		err := f.game.Start()
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func (f *Farming) Stop() error {
	f.game.Signal <- os.Interrupt

	err := f.game.Kill()
	if err != nil {
		return err
	}
	return nil
}
