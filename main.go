package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log.Println("Starting process")

	server, err := startServer("D:\\Jeux\\Farming Simulator 22\\dedicatedServer.exe")
	if err != nil {
		return
	}

	log.Println("Process running in background with PID:", server.Process.Pid)

	err = server.Wait()
	if err != nil {
		log.Println("Process crash... Restart process in 5 seconds...")
		time.Sleep(5 * time.Second)
		_ = clearTerminal()

		main()
	}

	log.Println("Process ended")
}

func clearTerminal() error {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Print("\033[H\033[2J") //clear terminal
		return err
	}

	return nil
}

func startServer(p string) (*exec.Cmd, error) {
	cmd, err := startProcess(p)
	if err != nil {
		log.Println(err)
	}

	newSignalHandler(func(s os.Signal) {
		cmd.Process.Kill()
	})

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, err
}

func startProcess(p string) (*exec.Cmd, error) {
	//run process in background
	cmd := exec.Command(p)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CmdLine:    "/c " + p,
	}
	return cmd, nil
}

func newSignalHandler(f func(os.Signal)) {
	go func(f func(os.Signal)) {
		for {
			c := make(chan os.Signal, 1)
			signal.Notify(c)

			f(<-c)
			break
		}
	}(f)
}
