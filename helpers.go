package main

import (
	"os"
	"os/user"
	"path/filepath"
)

func userFolder() string {
	user, err := user.Current()
	panicOnError(err)
	return user.HomeDir
}

func bellboyDirPath() string {
	return filepath.Join(userFolder(), ".bellboy")
}

func bellboyConfigPath() string {
	return filepath.Join(bellboyDirPath(), "bellboy.yaml")
}

func ensureBellboyDir() {
	_, errStatus := os.Stat(bellboyDirPath())
	if errStatus == nil {
		return
	}
	err := os.Mkdir(bellboyDirPath(), 0744)
	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
