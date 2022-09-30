package version

type Version interface {
	BinaryName() string
	Names() []string
	String() string
}

func FullName(v Version) string {
	return v.Names()[0]
}
