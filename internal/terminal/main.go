package terminal

import (
	"fmt"
	"os"
	"os/exec"
)

func Clean() error {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Print("\033[H\033[2J") //clear terminal
		return err
	}

	return nil
}

func Title(t string) {
	cmd := exec.Command("title", t)
	_ = cmd.Run()
}
