package server

import (
	"github.com/MatroxMC/FS22ServerManager/internal/event"
	"github.com/kataras/golog"
	"net/http"
	"net/url"
	"time"
)

type Api struct {
	server   *Server
	Password string
	Username string
	Ready    chan bool
	Logged   bool
	cookie   string
	handler  event.Handler
}

const (
	Url      = "http://192.168.1.42:8080"
	LoginURL = "/index.html"
)

func NewApi(s *Server) (*Api, error) {
	d := &Api{
		server:   s,
		Password: "secret",
		Username: "test",
		Ready:    make(chan bool),
	}

	return d, nil
}

func (a *Api) StartDaemon() {
	go func() {

		golog.Info("Starting API daemon")

		for {
			select {
			case <-a.server.Signal:
				golog.Debug("api task received signal")
				return

			default:

				s, err := a.HTTPStatus()
				if err != nil {
					continue
				}

				if s == 200 {
					a.Ready <- true
				}

				time.Sleep(1 * time.Second)
			}
		}

	}()
}

func (a *Api) HTTPStatus() (int, error) {
	resp, err := http.Get(Url)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func (a *Api) Login() error {

	u, err := url.JoinPath(Url, LoginURL)
	if err != nil {
		return err
	}

	resp, err := http.PostForm(u, url.Values{
		"username": {a.server.DedicatedConfig.Webserver.InitialAdmin.Username},
		"password": {a.server.DedicatedConfig.Webserver.InitialAdmin.Passphrase},
		"login":    {"Login"},
	})
	if err != nil {
		return err
	}

	c := resp.Cookies()
	for _, cookie := range c {
		if cookie.Name == "SessionID" {
			a.cookie = cookie.Value
			a.Logged = true

			golog.Debug("Logged in API")
			break
		}
	}

	return nil
}
