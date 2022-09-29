package version

type FS22 struct {
}

func (f FS22) String() string {
	return "FS22"
}

func (f FS22) BinaryName() string {
	return "dedicatedServer.exe"
}

func (f FS22) Names() []string {
	return []string{
		"Farming Simulator 22",
		"22",
		"FS22",
	}
}
