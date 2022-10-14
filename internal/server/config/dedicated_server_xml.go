package config

import (
	"encoding/xml"
	"io"
	"os"
	"path"
)

type DedicatedConfig struct {
	XMLName   xml.Name `xml:"server"`
	Text      string   `xml:",chardata"`
	Webserver struct {
		Text         string `xml:",chardata"`
		Port         string `xml:"port,attr"`
		InitialAdmin struct {
			Text       string `xml:",chardata"`
			Username   string `xml:"username"`
			Passphrase string `xml:"passphrase"`
		} `xml:"initial_admin"`
		Tls struct {
			Text                     string `xml:",chardata"`
			Port                     string `xml:"port,attr"`
			Active                   string `xml:"active,attr"`
			Certificate              string `xml:"certificate"`
			Privatekey               string `xml:"privatekey"`
			IntermediateCertificates struct {
				Text        string   `xml:",chardata"`
				Certificate []string `xml:"certificate"`
			} `xml:"intermediateCertificates"`
		} `xml:"tls"`
	} `xml:"webserver"`
	Game struct {
		Text        string `xml:",chardata"`
		Description string `xml:"description,attr"`
		Name        string `xml:"name,attr"`
		Exe         string `xml:"exe,attr"`
		Imprint     struct {
			Text   string `xml:",chardata"`
			Active string `xml:"active,attr"`
		} `xml:"imprint"`
		Logos struct {
			Text   string `xml:",chardata"`
			Login  string `xml:"login"`
			Bottom string `xml:"bottom"`
		} `xml:"logos"`
	} `xml:"game"`
}

func GetDedicatedConfig(directory string, file string) (DedicatedConfig, error) {
	x, err := os.Open(path.Join(directory, file))
	if err != nil {
		return DedicatedConfig{}, nil
	}

	b, err := io.ReadAll(x)
	if err != nil {
		return DedicatedConfig{}, nil
	}

	var p DedicatedConfig
	err = xml.Unmarshal(b, &p)
	if err != nil {
		return DedicatedConfig{}, nil
	}

	return p, nil
}
