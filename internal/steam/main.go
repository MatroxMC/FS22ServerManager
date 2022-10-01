package steam

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

type Steam bool

func (s Steam) ActiveUser() (string, error) {
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

func (s Steam) IsRunning() error {
	//thx to https://github.com/jshackles/idle_master/issues/217

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam\ActiveProcess`, registry.QUERY_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	pid, _, err := k.GetIntegerValue("pid")
	if err != nil {
		return err
	}

	login, _, err := k.GetIntegerValue("ActiveUser")
	if err != nil {
		return err
	}

	if pid == 0 || login == 0 {
		return fmt.Errorf("steam is not running")
	}

	return nil
}
