package farming

const ConfName = "config.toml"

type FarmingSimulator struct {
	Directory string `toml:"directory"`
	Steam     bool   `toml:"steam"`
	Version   string `toml:"version"`
}
