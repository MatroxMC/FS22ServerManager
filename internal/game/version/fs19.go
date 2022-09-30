package version

type FS19 struct{}

func (f FS19) BinaryName() string {
	return "dedicatedServer.exe"
}

func (f FS19) Names() []string {
	return []string{
		"Farming Simulator 19",
		"19",
		"FS19",
	}
}

func (f FS19) String() string {
	return "FS19"
}
