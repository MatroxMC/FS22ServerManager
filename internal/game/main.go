package game

import (
	"fmt"
	"github.com/MatroxMC/FS22ServerManager/internal/game/version"
	"github.com/MatroxMC/FS22ServerManager/internal/process"
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"log"
	"path"
)

type Version string
type Binary string
type Steam bool //TODO : ADD function to check if steams is available

// This var save all game version
var versions = []version.Version{
	version.FS22{},
	version.FS19{},
}

type Game struct {
	Binary     Binary
	Steam      Steam
	Version    version.Version
	Process    *process.Process
	ShowWindow bool
	Directory  string
	OnClose    func(Game, error) error
	OnStart    func(Game) error
}

func (g Game) Start() error {
	log.Print("Game starting...")
	p, err := process.New(g.Binary.String())
	if err != nil {
		return err
	}

	g.Process = p

	if g.OnStart != nil {
		err = g.OnStart(g)
		if err != nil {
			return err
		}
	}

	//Start the process
	err = g.Process.Start()
	if err != nil {
		return err
	}

	log.Print("Done.")

	err = g.Process.Wait()

	if g.OnClose != nil {
		return g.OnClose(g, err)
	}

	return nil
}

func (g Game) Restart() error {
	if g.Process.Running() {
		_ = g.Process.Kill()
	}

	return g.Start()
}

func New(directory string, version Version, steam Steam, window bool) (*Game, error) {
	//check if version exist
	v, err := version.Find()
	if err != nil {
		return nil, err
	}

	//check if game binary exist
	binary := Binary(path.Join(directory, v.BinaryName()))
	err = binary.Exist()
	if err != nil {
		return nil, err
	}

	//TODO : ADD MODS CHECK AND STEAM CHECK
	return &Game{
		Version:    v,
		Binary:     binary,
		Steam:      steam,
		Directory:  directory,
		ShowWindow: window,
	}, nil
}

func (v Version) Find() (version.Version, error) {
	for _, vv := range versions {
		for _, name := range vv.Names() {
			if name == string(v) {
				return vv, nil
			}
		}
	}

	return version.FS22{}, fmt.Errorf("version %s not found", v)
}

func (v Version) String() string {
	return string(v)
}

func (d Binary) Exist() error {
	return file.Exist(string(d))
}

func (d Binary) String() string {
	return string(d)
}
