package main

import (
	"github.com/BurntSushi/toml"
	"github.com/MatroxMC/FS22ServerManager/cmd/farming"
	"github.com/MatroxMC/FS22ServerManager/cmd/http"
	"github.com/MatroxMC/FS22ServerManager/internal/game"
	"github.com/MatroxMC/FS22ServerManager/internal/game/version"
	"github.com/MatroxMC/FS22ServerManager/internal/terminal"
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"log"
	"os"
)

var Config = Property{
	Game: farming.Farming{
		Steam:      true,
		Directory:  farming.GameDir,
		Version:    game.Version(version.FS22{}.String()),
		ShowWindow: true,
	},
	Web: http.Web{
		Port:     8080,
		Host:     "localhost",
		Password: "password",
	},
	Debug: 0,
}

type Property struct {
	Game  farming.Farming `toml:"game"`
	Web   http.Web        `toml:"web"`
	Debug int             `toml:"debug"`
}

func main() {
	err := SetLogFile("log.txt")
	if err != nil {
		log.Println("Error while setting log file : ", err)
	}

	//Init and load configuration file
	property, err := Config.Init()
	if err != nil {
		log.Fatal("Error while loading config file : ", err)
	}

	//set console name
	terminal.Title("FS22 Server Manager - " + property.Game.Version.String())

	//Start the web server
	_, err = property.Game.Start()
	if err != nil {
		log.Fatal("Error while starting the game : ", err)
	}
}

// Init function make or load the config file
func (p Property) Init() (*Property, error) {
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
