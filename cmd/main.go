package main

import (
	"github.com/MatroxMC/FS22ServerManager/cmd/farming"
	"github.com/MatroxMC/FS22ServerManager/cmd/http"
	"github.com/MatroxMC/FS22ServerManager/internal/game"
	"github.com/MatroxMC/FS22ServerManager/internal/game/version"
	"log"
	"sync"
)

var Game game.Game

var Config = game.Property{
	Game: farming.FarmingSimulator{
		Steam:     true,
		Directory: farming.GameDir,
		Version:   version.FS22{}.String(),
	},
	Web: http.Web{
		Port:     8080,
		Host:     "localhost",
		Password: "password",
	},
}

var Group = sync.WaitGroup{}

func main() {
	farming.Logi("Server Started")
	property, err := Config.Init()
	if err != nil {
		log.Fatal(err)
	}

	d, err := game.New(property)
	if err != nil {
		log.Fatalf("Error while loading game: %v", err)
	}

	Game = *d

	Group.Add(1)
	_, err = Game.Start(&Group)
	if err != nil {
		log.Println(err)
	}

	Group.Wait()

	log.Println("Server stopped")

	farming.Logi("Server stopped")
}
