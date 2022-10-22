package farming

import (
	"github.com/MatroxMC/FS22ServerManager/internal/event"
	"github.com/MatroxMC/FS22ServerManager/internal/server"
	"github.com/MatroxMC/FS22ServerManager/internal/steam"
	"github.com/kataras/golog"
	"sync"
)

/*
 * This Farming struct is used to store the configuration of the server in config file
 *
 * Directory: Is the directory where the server is installed
 * Steam : true if the server is installed on steam, if steam is not started the process wait steam to start
 * Window : true if the server window is shown
 */

type Farming struct {
	Directory   string      `toml:"directory"`
	Steam       steam.Steam `toml:"steam"`
	Window      bool        `toml:"window"`
	RestartTime int         `toml:"restart_time"`
	server      server.Server
}

func (f *Farming) Init() error {
	g, err := server.NewServer(f.Directory, server.DefaultInfo().DedicatedBinary, f.Steam, f.Window)
	if err != nil {
		return err
	}

	f.ProgramHandler().NewHandler(func() error {
		golog.Debug("TODO: Restart the server")
		return nil
	}, server.ProgramClosed)

	//Init default variables
	f.server = *g
	return nil
}

func (f *Farming) Start(w *sync.WaitGroup) error {
	err := f.Init()
	if err != nil {
		return err
	}

	//Run all program in computer
	err = f.server.Run()
	if err != nil {
		return err
	}

	err = f.server.Start()
	if err != nil {
		return err
	}

	w.Add(1)
	go func() {
		defer w.Done()
		golog.Debug("server running in goroutine")

		//Wait for the server to stop
		err = f.server.Wait()
		if err != nil {
			golog.Error(err)
		}
	}()

	return nil
}

func (f *Farming) Stop() error {
	return f.server.Stop()
}

func (f *Farming) ServerHandler() *event.Handler {
	return &f.server.Handler
}

func (f *Farming) ProgramHandler() *event.Handler {
	return &f.server.Program.Handler
}
