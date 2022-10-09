package main

import (
	"github.com/MatroxMC/FS22ServerManager/cmd/farming"
	"github.com/MatroxMC/FS22ServerManager/cmd/http"
	"github.com/MatroxMC/FS22ServerManager/internal/terminal"
	"github.com/kataras/golog"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	ConfName = "config.toml"
)

var service = &Config{
	Farming: &farming.Farming{
		Directory:   "D:\\Jeux\\Farming Simulator 22",
		Steam:       true,
		Window:      true,
		RestartTime: 5,
	},
	Http: &http.Http{
		Address: "127.0.0.1",
		Port:    8080,
	},
	Log: Log{
		Level: "info",
	},
}

var waitGroup = &sync.WaitGroup{}

func main() {
	handleSignal()
	_, _ = terminal.Title("Farming Simulator Server Manager")
	err := initConfig()
	if err != nil {
		panic(err)
	}

	golog.SetLevel(service.Log.Level)

	err = service.Http.Start(waitGroup)
	if err != nil {
		panic(err)
	}

	err = service.Farming.Start(waitGroup)
	if err != nil {
		panic(err)
	}

	waitGroup.Wait()
}

func handleSignal() {
	c := make(chan os.Signal, 1)

	signal.Notify(c)

	go func() {
		_ = <-c

		golog.Info("Stop all services...")

		go func() {
			err := service.Farming.Stop()
			if err != nil {
				golog.Print(err)
			}

			golog.Debug("Server Manager stopped")
		}()

		go func() {
			err := service.Http.Stop()
			if err != nil {
				golog.Print(err)
			}

			golog.Debug("Mod API stopped")
		}()

		time.AfterFunc(time.Second*5, func() {
			golog.Warn("Force stop")
			os.Exit(0)
		})
	}()
}
