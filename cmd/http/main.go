package http

import (
	"crypto/md5"
	"encoding/hex"
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
		f := "C:\\Users\\matro\\Documents\\My Games\\FarmingSimulator2022\\mods"

		golog.Debug("getting mods in %s", f)

		ParseMods(f)

		http.HandleFunc("/", handleHttp)
		serve := http.Server{
			Addr: h.Address + ":80",
		}

		//Register dir access
		fs := http.FileServer(http.Dir(f))
		http.Handle("/mods/", http.StripPrefix("/mods", fs))

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

	golog.Debug("mods Manager started")

	return nil
}

func (h *Http) Stop() error {

	if h.serve != nil {
		golog.Debug("stop http server")
		err := h.serve.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseMods(s string) {
	golog.Debug("parsing mod(s)")
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

		n := s + "\\" + f.Name()

		st, err := os.Stat(n)
		if err != nil {
			log.Print(err)
			continue
		}

		hash, err := md5sum(f.Name())
		if err != nil {
			log.Print(err)
			continue
		}

		files = append(files, Mod{
			Name: f.Name(),
			Path: "http://127.0.0.1/mods/" + f.Name(),
			Size: st.Size(),
			Hash: hash,
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

func md5sum(p string) (string, error) {
	hash := md5.New()
	hash.Write([]byte(p))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
