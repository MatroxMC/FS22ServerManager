package http

type Web struct {
	Port     int    `toml:"port"`
	Host     string `toml:"host"`
	Password string `toml:"password"`
}
