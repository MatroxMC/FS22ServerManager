package main

import (
	"github.com/BurntSushi/toml"
	"github.com/MatroxMC/FS22ServerManager/cmd/farming"
	"github.com/MatroxMC/FS22ServerManager/cmd/http"
	"github.com/MatroxMC/FS22ServerManager/internal/server"
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

type Log struct {
	Level string `toml:"level"`
}

type Config struct {
	Farming *farming.Farming `toml:"farming"`
	Http    *http.Http       `toml:"http"`
	Api     *server.Api      `toml:"api"`
	Log     Log              `toml:"log"`
}

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

// Init function make or load the config file
func initConfig() error {

	//If the config file does not exist, create it
	if _, err := os.Stat(ConfName); os.IsNotExist(err) {
		f, err := os.Create(ConfName)
		defer f.Close()
		if err != nil {
			return err
		}

		//write default config
		err = toml.NewEncoder(f).Encode(service)
		if err != nil {
			return err
		}

		return nil
	}

	_, err := toml.DecodeFile(ConfName, service)
	if err != nil {
		return err
	}

	return nil
}

func handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	go func() {
		_ = <-c

		golog.Debug("signal received")

		service.Farming.Stop()
		service.Http.Stop()

		time.AfterFunc(time.Second*5, func() {
			golog.Warn("FORCE CLOSE")
			os.Exit(0)
		})
	}()
}
