package main

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/work"
	"github.com/fpawel/goutils/panichook"
	"github.com/fpawel/goutils/winapp"
	"github.com/lxn/win"
	"os"
	"path/filepath"
)

func main() {

	// Преверяем, не было ли приложение запущено ранее
	hWndPeer := work.FindPeer()
	if winapp.IsWindow(hWndPeer) {
		// Если было, выдвигаем окно приложения на передний план и завершаем процесс
		win.ShowWindow(hWndPeer, win.SW_RESTORE)
		win.SetForegroundWindow(hWndPeer)
		return
	}

	const serverAppExe = "anbus.exe"
	dir := filepath.Dir(os.Args[0])
	if _, err := os.Stat(filepath.Join(dir, serverAppExe)); os.IsNotExist(err) {
		dir = anbus.AppName.Dir()
	}
	panichook.Run(filepath.Join(dir, serverAppExe))
}
