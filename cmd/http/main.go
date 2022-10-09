package http

import (
	"encoding/json"
	"github.com/kataras/golog"
	"log"
	"net/http"
	"os"
	"sync"
)

type Http struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	serve   *http.Server
}

type Mod struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
}

var mods []Mod

func (h *Http) Start(w *sync.WaitGroup) error {
	w.Add(1)
	go func() {

		ParseMods("C:\\Users\\matro\\Documents\\My Games\\FarmingSimulator2022\\mods")

		http.HandleFunc("/", handleHttp)
		serve := http.Server{
			Addr: h.Address + ":80",
		}

		h.serve = &serve
		err := serve.ListenAndServe()
		if err != nil {
			switch err {
			case http.ErrServerClosed:
				break
			default:
				panic(err)
			}
		}

		w.Done()
	}()

	golog.Info("Mods Manager started")

	return nil
}

func (h *Http) Stop() error {
	err := h.serve.Close()
	if err != nil {
		return err
	}
	return nil
}

func ParseMods(s string) {
	dir, err := os.ReadDir(s)
	if err != nil {
		log.Print(err)
		return
	}

	var files = make([]Mod, 0)
	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		files = append(files, Mod{
			Name: f.Name(),
			Path: s + f.Name(),
			Size: 0,
		})
	}

	golog.Info("Parsed mods: ", len(files))

	mods = files
}

func handleHttp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(mods)
	if err != nil {
		return
	}
}
