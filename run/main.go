package main

import (
	"bytes"
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/goutils/panichook"
	"github.com/fpawel/goutils/winapp"
	"github.com/lxn/win"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	dir := filepath.Dir(os.Args[0])
	const (
		serverAppExe = "anbus.exe"
		peerAppExe   = "anbusui.exe"
	)

	if _, err := os.Stat(filepath.Join(dir, serverAppExe)); os.IsNotExist(err) {
		dir = anbus.AppName.Dir()
	}
	go panichook.Run(filepath.Join(dir, serverAppExe))

	if err := runUIApp(); err != nil {
		winapp.MsgBox(err.Error(), filepath.Join(dir, peerAppExe), win.MB_ICONERROR)
	}
	closeServerApp()
}

func closeServerApp() {
	win.SendMessage(winapp.FindWindow(anbus.ServerAppWindow), win.WM_CLOSE, 0, 0)
}

func runUIApp() error {
	cmd := exec.Command(anbus.AppName.FileName("anbusui.exe"), "-wait-server")
	cmd.Stdout = os.Stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
