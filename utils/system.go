package utils

import (
	"fmt"
	"os"
)

func GetSystemPath() string {
	home, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/%s", home, "kickshaw/")
}

func hasAppDir() bool {
	_, err := os.Stat(GetSystemPath())
	return !os.IsNotExist(err)
}

func mkdir() error {
	return os.Mkdir(GetSystemPath(), 0644)
}

func HasSystemPath() {
	if !hasAppDir() {
		if err := mkdir(); err != nil {
			HandleError(err)
		}
	}
}
