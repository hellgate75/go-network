package io

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Gets the Shared libraries extension included by dot, related to current O/S
func GetShareLibExt() string {
	if runtime.GOOS == "windows" {
		return ".dll"
	}
	return ".so"
}

// Retrieves Current Path
func GetCurrentPath() string {
	wd, err := os.Getwd()
	if err != nil {
		exec, err := os.Executable()
		if err != nil {
			return HomeFolder()
		}
		return filepath.Dir(exec)
	}
	return wd
}

// Retrieves Home directory
func HomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		return os.TempDir()
	}
	return usr.HomeDir
}

func TempFolder() string {
	return os.TempDir()
}

func UniqueTempFolder(suffix string) string {
	dir, _ := ioutil.TempDir(TempFolder(), suffix)
	return dir
}
