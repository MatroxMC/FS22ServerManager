package server

import (
	"fmt"
	"github.com/MatroxMC/FS22ServerManager/internal/event"
	"github.com/MatroxMC/FS22ServerManager/internal/server/config"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"github.com/kataras/golog"
	"os"
	"time"
)

type Server struct {
	Program         Program
	DedicatedConfig config.DedicatedConfig
	GameConfig      config.GameConfig
	Api             Api
	Info            Info
	Handler         event.Handler
	Signal          chan os.Signal
}

type Info struct {
	DedicatedBinary string
	Names           []string
	String          string
	DataPath        string
	DedicatedConfig string
	GameConfig      string
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

	e, err := config.GetGameConfig(DefaultInfo().DataPath, DefaultInfo().GameConfig)
	if err != nil {
		return &Server{}, err
	}

	s := &Server{
		DedicatedConfig: d,
		GameConfig:      e,
		Program:         p,
		Info:            DefaultInfo(),
		Signal:          make(chan os.Signal, 1),
	}

	api, err := NewApi(s)
	if err != nil {
		return &Server{}, err
	}

	s.Api = *api

	//Return a clean Server struct
	return s, nil
}

func (s *Server) Run() error {
	err := s.Program.Start()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Start() error {
	st := s.Program.Steam

	if st {
		if !st.IsRunning() {
			golog.Info("steam is not running, waiting for steam to start")
			r := st.WaitForRunning()

			if r == steam.StatusExited {
				return nil
			}

			golog.Info("steam is running")
		}
	}

	for status := s.IsStarted(); !status; {
		select {
		case <-s.Signal:
			return nil
		default:
			time.Sleep(1 * time.Second)
			status = s.IsStarted()
		}
	}

	golog.Debug("api is running")

	err := s.Api.Login()
	if err != nil {
		return err
	}

	golog.Debug("logged in to api cookie: ", s.Api.cookie)

	go func() {
		for true {
			select {
			case <-s.Signal:
				return
			default:
				stats := s.Api.IsOnline()

				if !stats {

				}
			}

			time.Sleep(1 * time.Second)
		}

		golog.Info("stop go routine status checker")
	}()

	err = s.Api.Start()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) IsStarted() bool {
	_, err := s.Api.HTTPStatus()
	if err != nil {
		return false
	}

	return true
}

func (s *Server) Stop() error {
	s.Signal <- os.Interrupt
	return s.Program.Stop()
}

func (s *Server) Wait() error {
	return s.Program.Cmd.Wait()
}

func (s *Server) Restart() {
	golog.Info("restarting server")
}

func DefaultInfo() Info {
	return Info{
		DedicatedBinary: "dedicatedServer.exe",
		DedicatedConfig: "dedicatedServer.xml",
		GameConfig:      "gameSettings.xml",
		Names: []string{
			"Farming Simulator 22",
			"22",
			"FS22",
		},
		String:   "Farming Simulator 22",
		DataPath: "C:\\Users\\%username%\\Documents\\My Games\\FarmingSimulator2022",
	}
}
