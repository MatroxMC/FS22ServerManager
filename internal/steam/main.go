package steam

import (
	"golang.org/x/sys/windows/registry"
	"os"
)

type Steam bool

func (s Steam) GetAutoLoginUser() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
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
	//thx to https://github.com/jshackles/idle_master/issues/217

	pid, err := s.GetPID()
	if err != nil {
		return false
	}

	//if pid is 0, steam is not running
	if pid == 0 {
		return false
	}

	//check if process is running
	_, err = os.FindProcess(int(pid))
	if err != nil {
		return false
	}

	return true
}

func (s Steam) GetPID() (int, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam\ActiveProcess`, registry.QUERY_VALUE)
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
