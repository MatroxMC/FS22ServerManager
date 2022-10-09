package main

import (
	"github.com/MatroxMC/FS22ServerManager/internal/game"
	"testing"
	"time"
)

func TestCreateGame(t *testing.T) {
	g, err := game.New("D:\\Jeux\\Farming Simulator 22", false, true)
	if err != nil {
		return
	}

	g.NewHandler(func() error {
		return nil
	}, game.HandleStart)

	g.NewHandler(func() error {
		return nil
	}, game.HandleStop)

	g.NewHandler(func() error {
		return nil
	}, game.HandleClosed)

	go func() {
		time.Sleep(5 * time.Second)
		err := g.Stop()
		if err != nil {
			t.Error(err)
		}
	}()

	err = g.Start()
	if err != nil {
		t.Error(err)
	}
}
