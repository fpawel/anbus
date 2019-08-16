package app

import (
	"context"
	"github.com/fpawel/anbus/internal"
	"github.com/fpawel/anbus/internal/peer"
	"github.com/fpawel/comm/comport"
	"github.com/fpawel/dseries"
	"github.com/lxn/win"
	"github.com/powerman/structlog"
	"path/filepath"
	"sync"
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

var (
	port           *comport.ReadWriter
	chRequest      chan modbusRequest
	ctxApp         context.Context
	cancelWorkFunc = func() {}
	skipDelayFunc  = func() {}
	wgWork         sync.WaitGroup
	log            = structlog.New()
)
