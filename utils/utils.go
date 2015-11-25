package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

const PROGRAM_NAME = "helm"

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func GetConfigPath() string {
	return filepath.Join(getHomeDir(), "."+PROGRAM_NAME)
}
