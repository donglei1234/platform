package fxapp

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/fx"
)

type AppName struct {
	fx.In
	AppName string `name:"AppName"`
}

type AppNameLoader struct {
	fx.Out
	AppName string `name:"AppName"`
}

func (l *AppNameLoader) LoadConstant(value string) error {
	l.AppName = value
	return nil
}

func (l *AppNameLoader) LoadFromExecutable() (err error) {
	if exeName, e := os.Executable(); e != nil {
		err = e
	} else {
		exeName = filepath.Base(exeName)

		if runtime.GOOS == "windows" && strings.HasSuffix(exeName, ".exe") {
			exeName = exeName[:len(exeName)-4]
		}

		l.AppName = exeName
	}
	return
}
