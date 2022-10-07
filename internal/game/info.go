package game

type Info struct {
	Binary string
	Names  []string
	String string
}

func DefaultInfo() Info {
	return Info{
		Binary: "dedicatedServer.exe",
		Names: []string{
			"Farming Simulator 22",
			"22",
			"FS22",
		},
		String: "Farming Simulator 22",
	}
}
