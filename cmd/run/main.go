package main

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/goutils/panichook"
	"os"
	"path/filepath"
)

func main() {
	const serverAppExe = "anbus.exe"
	dir := filepath.Dir(os.Args[0])
	if _, err := os.Stat(filepath.Join(dir, serverAppExe)); os.IsNotExist(err) {
		dir = anbus.AppName.Dir()
	}
	panichook.Run(filepath.Join(dir, serverAppExe))
}
