package main

import (
	"github.com/BurntSushi/toml"
	"github.com/MatroxMC/FS22ServerManager/cmd/farming"
	"github.com/MatroxMC/FS22ServerManager/cmd/http"
	"github.com/MatroxMC/FS22ServerManager/internal/server"
	"os"
)

type Config struct {
	Farming *farming.Farming `toml:"farming"`
	Http    *http.Http       `toml:"http"`
	Api     *server.Api      `toml:"api"`
	Log     Log              `toml:"log"`
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
