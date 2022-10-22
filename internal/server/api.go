package server

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

type Api struct {
	server *Server
	Logged bool
	cookie string
}

const (
	Url      = "http://192.168.1.42:8080"
	IndexURL = "/index.html"
)

func NewApi(s *Server) (*Api, error) {
	d := &Api{
		server: s,
	}

	return d, nil
}

func (a *Api) HTTPStatus() (int, error) {
	resp, err := http.Get(Url)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func (a *Api) Start() error {
	u, err := url.JoinPath(Url, IndexURL)
	if err != nil {
		return err
	}

	c := &http.Client{}
	req, err := http.NewRequest("POST", u, strings.NewReader(url.Values{
		"game_name":           {"Managed By FS22ServerManager"},
		"admin_password":      {a.server.DedicatedConfig.Webserver.InitialAdmin.Passphrase},
		"game_password":       {a.server.GameConfig.CreateGame.Password},
		"savegame":            {a.server.GameConfig.CreateGame.Name},
		"map_start":           {"default_MapUS"},
		"difficulty":          {"2"},
		"server_port":         {a.server.GameConfig.CreateGame.Port},
		"max_player":          {a.server.GameConfig.CreateGame.Capacity},
		"mp_language":         {a.server.GameConfig.MpLanguage},
		"auto_save_interval":  {"180"},
		"stats_interval":      {"360"},
		"pause_game_if_empty": {"2"},
		"crossplay_allowed":   {a.server.GameConfig.CreateGame.AllowCrossPlay},
		"start_server":        {"Start"},
	}.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "SessionID="+a.cookie)

	//send request
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (a *Api) Login() error {

	u, err := url.JoinPath(Url, IndexURL)
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
			break
		}
	}

	return nil
}

func (a *Api) IsOnline() bool {
	h, err := http.Get(Url)
	if err != nil {
		return false
	}
	if h.StatusCode != 200 {
		return false
	}

	doc, err := goquery.NewDocumentFromReader(h.Body)
	if err != nil {
		return false
	}

	o := false
	doc.Find(".status-indicator").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "ONLINE" {
			o = true
		}
	})

	return o
}
