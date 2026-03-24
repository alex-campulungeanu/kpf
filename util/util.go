package util

import "os"

func HomeDir() (string, error) {
	return os.UserHomeDir()
}
