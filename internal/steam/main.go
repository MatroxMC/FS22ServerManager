package steam

import (
	"golang.org/x/sys/windows/registry"
	"os"
	"os/signal"
	"syscall"
)

type Status int

const (
	StatusOK     Status = 0x0
	StatusExited Status = 0x1
	RegistryKey         = registry.CURRENT_USER
	RegistryPath        = `Software\Valve\Steam`
)

type Steam bool

func (s Steam) WaitForRunning() Status {
	exit := make(chan os.Signal, 1)

	go func() {
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	}()

	for {
		select {
		case <-exit:
			signal.Stop(exit)
			close(exit)
			return StatusExited
		default:
			if s.IsRunning() {
				return StatusOK
			}
		}
	}
}

func (s Steam) GetExe() (string, error) {
	k, err := registry.OpenKey(RegistryKey, RegistryPath, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	p, _, err := k.GetStringValue("SteamExe")

	return p, err
}

func (s Steam) IsInstalled() bool {
	p, err := s.GetExe()
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return true
}

func (s Steam) GetAutoLoginUser() (string, error) {
	k, err := registry.OpenKey(RegistryKey, RegistryPath, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	u, _, err := k.GetStringValue("AutoLoginUser")
	if err != nil {
		return "", err
	}

	return u, nil
}

func (s Steam) IsRunning() bool {
	//https://github.com/jshackles/idle_master/issues/217
	pid, err := s.GetPID()
	if err != nil {
		return false
	}

	//if pid is 0, steam is not running
	if pid == 0 {
		return false
	}

	//check if process is running
	_, err = os.FindProcess(pid)
	if err != nil {
		return false
	}

	return true
}

func (s Steam) GetPID() (int, error) {

	k, err := registry.OpenKey(RegistryKey, RegistryPath+`\ActiveProcess`, registry.QUERY_VALUE)
	if err != nil {
		return 0, err
	}
	defer k.Close()

	//get pid process in registry
	pid, _, err := k.GetIntegerValue("pid")
	if err != nil {
		return 0, err
	}

	return int(pid), nil
}
