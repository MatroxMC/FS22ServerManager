package dir

import "os"

func Exist(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return err
	}
	return nil
}
