package app

import (
	"context"
	"github.com/fpawel/anbus/internal"
	"github.com/fpawel/anbus/internal/api/notify"
	"github.com/fpawel/anbus/internal/cfg"
	"github.com/fpawel/comm"
	"github.com/fpawel/comm/comport"
	"github.com/fpawel/dseries"
	"github.com/fpawel/gohelp/winapp"
	"github.com/hako/durafmt"
	"github.com/lxn/win"
	"github.com/powerman/structlog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Run() {

	// Преверяем, не было ли приложение запущено ранее.
	// Если было, выдвигаем окно UI приложения на передний план и завершаем процесс.
	if notify.ServerWindowAlreadyExists {
		hWnd := winapp.FindWindow(notify.PeerWindowClassName)
		win.ShowWindow(hWnd, win.SW_RESTORE)
		win.SetForegroundWindow(hWnd)
		log.Fatal("mil82.exe already executing")
	}

	dseries.Open(filepath.Join(internal.DataDir(), "series.sqlite"))
	log.Info("updated at", "time", dseries.UpdatedAt(), "elapsed", durafmt.Parse(time.Since(dseries.UpdatedAt())))

	var cancel func()
	ctxApp, cancel = context.WithCancel(context.TODO())
	closeHttpServer := startHttpServer()

	if os.Getenv("ANBUS_SKIP_RUN_PEER") != "true" {
		if err := exec.Command(filepath.Join(filepath.Dir(os.Args[0]), "anbusgui.exe")).Start(); err != nil {
			panic(err)
		}
	}

	// оснавная работа
	go work()

	// цикл оконных сообщений
	for {
		var msg win.MSG
		if win.GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
	cancel()
	closeHttpServer()
	notify.Window.Close()
	log.ErrIfFail(dseries.Close)
}

var (
	comPort = comport.NewReadWriter(func() comport.Config {
		c := cfg.Get()
		return comport.Config{
			Baud:        c.ComportBaud,
			Name:        c.ComportName,
			ReadTimeout: time.Millisecond,
		}
	}, func() comm.Config {
		return cfg.Get().Comm
	})

	chTasks         = make(chan func(), 1000)
	chConfigChanged = make(chan struct{})
	ctxApp          context.Context
	log             = structlog.New()
)
