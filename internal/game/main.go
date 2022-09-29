package game

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/MatroxMC/FS22ServerManager/cmd/farming"
	"github.com/MatroxMC/FS22ServerManager/cmd/http"
	"github.com/MatroxMC/FS22ServerManager/internal/game/version"
	"github.com/MatroxMC/FS22ServerManager/internal/process"
	"github.com/MatroxMC/FS22ServerManager/internal/terminal"
	"github.com/MatroxMC/FS22ServerManager/internal/tools/dir"
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
)

type Game struct {
	Version  version.Version
	Path     string
	Steam    bool
	Property *Property
}

func (g Game) Start(w *sync.WaitGroup) (*process.Process, error) {
	farming.Logi("Start process")

	log.Println("Starting Web Server")

	bin := filepath.Join(g.Path, g.Version.BinaryName())
	pp, err := process.NewProcess(bin)
	if err != nil {
		return nil, err
	}

	c := make(chan os.Signal, 1)

	//This routine is used to check the signal and kill the process
	go func() {
		for {
			signal.Notify(c)

			fmt.Println(c)

			switch <-c {
			case os.Interrupt:
			case os.Kill:

				err := pp.Stop()
				if err != nil {
					fmt.Println("Error while stopping process: ", err)
				}
			}

			pp.Cmd.Process.Kill()

			farming.Logi("End signal routine")
			break
		}
	}()

	w.Add(1)
	go func() {
		err := pp.Start()
		if err != nil {
			w.Done()
			return
		}

		log.Println("Web Server running")

		err = pp.Cmd.Wait()
		if err != nil {
			log.Println("Server stopped, restarting in 2 seconds...")

			_ = terminal.Clean()
			_, err = g.Start(w) //restart server
			if err != nil {
				log.Fatal("Error while restarting server: ", err)
			}
		}

		w.Done()

		signal.Stop(c) //stop listening for signals if process is stopped
		close(c)

		farming.Logi("Close signal routine 2")
	}()

	return pp, nil
}

type Property struct {
	Game farming.FarmingSimulator `toml:"game"`
	Web  http.Web                 `toml:"web"`
}

func (p Property) Init() (*Property, error) {

	//create game config if not exist
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

func New(property *Property) (*Game, error) {

	v, err := version.FindByString(property.Game.Version)
	if err != nil {
		return nil, err
	}

	//check if game dir exist
	err = dir.Exist(property.Game.Directory)
	if err != nil {
		return nil, err
	}

	//check if mods dir exist
	err = dir.Exist(property.Game.Directory)
	if err != nil {
		log.Println("!!WARNING!! Mods directory not found")
	}

	//check if game binary exist
	serverBinary := filepath.Join(property.Game.Directory, v.BinaryName())
	err = file.Exist(serverBinary)
	if err != nil {
		return nil, err
	}

	return &Game{
		Version:  v,
		Path:     property.Game.Directory,
		Steam:    property.Game.Steam,
		Property: property,
	}, nil
}
