package steam

import (
	"github.com/MatroxMC/FS22ServerManager/internal/tools/file"
	"golang.org/x/sys/windows/registry"
	"os"
	"time"
)

const (
	RegistryKey  = registry.CURRENT_USER
	RegistryPath = `Software\Valve\Steam`
)

type Steam bool

func (s Steam) Wait() error {

	var run = false //If the game is running
	var msg = false //If the message has been sent

	for !run {
		run = s.IsRunning()
		if !run {
			if !msg {
				msg = true
			}
			time.Sleep(time.Second * 2)
		}
	}

	return nil
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
	err = file.Exist(p)
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
