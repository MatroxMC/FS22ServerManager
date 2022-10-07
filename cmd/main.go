package main

import (
	"github.com/BurntSushi/toml"
	"github.com/MatroxMC/FS22ServerManager/cmd/farming"
	"github.com/MatroxMC/FS22ServerManager/internal/terminal"
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"log"
	"os"
	"sync"
)

var Config = Property{
	Game: farming.Farming{
		Steam:     true,
		Directory: farming.GameDir,
		Window:    true,
	},
	Debug: 0,
}

type Property struct {
	Game  farming.Farming `toml:"game"`
	Debug int             `toml:"debug"`
}

var waitGroup = sync.WaitGroup{}

func main() {
	err := SetLogFile("debug.log")
	if err != nil {
		log.Println("Error while setting log file : ", err)
	}

	//Init and load configuration file
	property, err := Config.init()
	if err != nil {
		log.Fatal("Error while loading config file : ", err)
	}

	waitGroup.Add(1)
	go func() {
		//Run the web server
		game, err := property.Game.Start()
		if err != nil {
			log.Fatal("Error while starting the game : ", err)
		}

		//set console name
		_, _ = terminal.Title(game.Info.String + " Server Manager")

		game.Process.Wait()

		waitGroup.Done()
	}()

	waitGroup.Wait()
}

// Init function make or load the config file
func (p Property) init() (*Property, error) {
	err := file.Exist(farming.ConfName)
	if err != nil {
		f, err := os.Create(farming.ConfName)
		defer f.Close()
		if err != nil {
			return nil, err
		}

		//write default config
		err = toml.NewEncoder(f).Encode(p)
		if err != nil {
			return nil, err
		}

		return &p, nil
	}

	var config Property
	_, err = toml.DecodeFile(farming.ConfName, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
