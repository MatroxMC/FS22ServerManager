package http

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Http struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

func (h Http) Start(w *sync.WaitGroup) error {

	w.Add(1)
	go func() {

		http.HandleFunc("/", handleHttp)
		err := http.ListenAndServe(h.Address+":80", nil)
		if err != nil {
			panic(err)
		}

		w.Done()
	}()

	return nil
}

func handleHttp(w http.ResponseWriter, req *http.Request) {
	j := map[string]string{
		"message": "Hello World",
		"status":  "ok",
		"code":    "200",
		"method":  req.Method,
		"url":     req.URL.String(),
		"host":    req.Host,
		"proto":   req.Proto,
		"remote":  req.RemoteAddr,
		"agent":   req.UserAgent(),
		"referer": req.Referer(),
		"header":  req.Header.Get("Content-Type"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(j)
	if err != nil {
		return
	}
}
