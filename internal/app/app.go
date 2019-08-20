package app

import (
	"context"
	"github.com/fpawel/anbus/internal"
	"github.com/fpawel/anbus/internal/cfg"
	"github.com/fpawel/anbus/internal/peer"
	"github.com/fpawel/comm"
	"github.com/fpawel/comm/comport"
	"github.com/fpawel/dseries"
	"github.com/lxn/win"
	"github.com/powerman/structlog"
	"path/filepath"
	"time"
)

const (
	pipeName                      = `\\.\pipe\anbus`
	anbusServerAppWindowClassName = "AnbusServerAppWindow"
	peerWindowClassName           = "TAnbusMainForm"
)

func Run() {

	peer.AssertRunOnce()

	dseries.Open(filepath.Join(internal.DataDir(), "mil82.series.sqlite"))

	var cancel func()
	ctxApp, cancel = context.WithCancel(context.TODO())
	closeHttpServer := startHttpServer()
	peer.Init("")

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
	peer.Close()
	log.ErrIfFail(dseries.Close)
}

type peerNotifier struct{}

func (_ peerNotifier) OnStarted() {
	peer.InitPeer()
}

func (_ peerNotifier) OnClosed() {
	peer.ResetPeer()
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

	chTasks = make(chan func(), 1000)
	ctxApp  context.Context
	log     = structlog.New()
)
