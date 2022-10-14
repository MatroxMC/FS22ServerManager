package server

import (
	"fmt"
	"github.com/MatroxMC/FS22ServerManager/internal/event"
	"github.com/MatroxMC/FS22ServerManager/internal/server/config"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"github.com/kataras/golog"
	"os"
)

type Server struct {
	Program         Program
	DedicatedConfig config.DedicatedConfig
	Api             Api
	Info            Info
	Handler         event.Handler
	Signal          chan os.Signal
	Killed          bool
}

type Info struct {
	DedicatedBinary string
	Names           []string
	String          string
	DataPath        string
	DedicatedConfig string
}

func NewServer(directory string, executable string, steam steam.Steam, window bool) (*Server, error) {
	//Check if the directory exist
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return &Server{}, fmt.Errorf("directory %s does not exist", directory)
	}

	p, err := NewProgram(directory, executable, steam, window)
	if err != nil {
		return &Server{}, err
	}

	d, err := config.GetDedicatedConfig(directory, DefaultInfo().DedicatedConfig)
	if err != nil {
		return &Server{}, err
	}

	//Return a clean Server struct
	return &Server{
		DedicatedConfig: d,
		Program:         p,
		Info:            DefaultInfo(),
		Signal:          make(chan os.Signal, 1),
	}, nil
}

func (s *Server) Init() error {

	return nil
}

func (s *Server) Run() error {
	err := s.Init()
	if err != nil {
		return err
	}

	err = s.Program.Start()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Start() error {
	ss := s.Program.Steam

	if ss {
		if !ss.IsRunning() {
			golog.Info("Steam is not running, waiting for steam to start")
			s := ss.WaitForRunning()

			if s == steam.StatusExited {
				return nil
			}
			golog.Info("Steam is running")
		}
	}

	return nil
}

func (s *Server) Stop() error {
	return s.Program.Stop()
}

func (s *Server) Wait() error {
	return s.Program.Cmd.Wait()
}

func DefaultInfo() Info {
	return Info{
		DedicatedBinary: "dedicatedServer.exe",
		DedicatedConfig: "dedicatedServer.xml",
		Names: []string{
			"Farming Simulator 22",
			"22",
			"FS22",
		},
		String:   "Farming Simulator 22",
		DataPath: "C:\\Users\\%username%\\Documents\\My Games\\FarmingSimulator2022",
	}
}
