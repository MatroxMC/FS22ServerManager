package version

import "fmt"

type Version interface {
	Names() []string
	String() string
	BinaryName() string
}

var Games = []Version{
	FS22{},
	FS19{},
}

func FindByString(s string) (Version, error) {
	for _, v := range Games {
		for _, n := range v.Names() {
			if s == n {
				return v, nil
			}
		}
	}
	return nil, fmt.Errorf("version %s not found", s)
}
