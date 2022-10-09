package main

import (
	"github.com/MatroxMC/FS22ServerManager/cmd/farming"
	"github.com/MatroxMC/FS22ServerManager/cmd/http"
	"log"
	"os"
	"os/signal"
	"sync"
)

const (
	ConfName = "config.toml"
)

var service = &Config{
	Farming: farming.Farming{
		Directory:   "D:\\Jeux\\Farming Simulator 22",
		Steam:       true,
		Window:      true,
		RestartTime: 5,
	},
	Http: http.Http{
		Address: "127.0.0.1",
		Port:    8080,
	},
}

var waitGroup = &sync.WaitGroup{}

func main() {
	handleSignal()

	err := initConfig()
	if err != nil {
		panic(err)
	}

	err = service.Farming.Start(waitGroup)
	if err != nil {
		panic(err)
	}
	log.Print("Farming Simulator Manager started")

	err = service.Http.Start(waitGroup)
	if err != nil {
		panic(err)
	}
	log.Print("Mod API started")

	waitGroup.Wait()
}

func handleSignal() {
	c := make(chan os.Signal, 1)

	signal.Notify(c)

	go func() {
		_ = <-c

		service.Farming.Stop()

		log.Print("Closing...")

	}()
}
