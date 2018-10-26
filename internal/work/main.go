package work

import (
	"fmt"
	"github.com/Microsoft/go-winio"
	"github.com/fpawel/goutils/serial/comport"
	"github.com/fpawel/anbus/internal/data"
	"github.com/fpawel/anbus/internal/notify"
	"github.com/fpawel/anbus/internal/panalib"
	"github.com/lxn/win"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"net"
	"sync"
)

func Main() {
	cfg, errCfg := panalib.OpenConfig()
	if errCfg != nil {
		fmt.Println(errCfg)
		cfg = panalib.DefaultConfig()
	}
	x := &worker{
		config:          cfg,
		chModbusRequest: make(chan modbusRequest, 10),
		series:          data.NewSeries(),
		comport:         comport.NewPortWithConfig(cfg.Comport),
		ln:              mustPipeListener(),
	}
	x.window = notify.NewWindow(x.onCommand)
	x.initPeer()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			switch conn, err := x.ln.Accept(); err {
			case nil:
				go jsonrpc2.ServeConn(conn)
			case winio.ErrPipeListenerClosed:
				return
			default:
				panic(err)
			}
		}
	}()

	go func() {
		x.mainWork()
		wg.Done()
	}()

	// цикл оконных сообщений
	for {
		var msg win.MSG
		if win.GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}

	// всё закрыть

	x.flagClose = true // установить флаг, сигнализирующий что надо выйти из всех бесконечных циклов

	if err := x.ln.Close(); err != nil {
		fmt.Println("close pipe listener error:", err)
	}

	x.comport.Cancel() // прервать СОМ порт
	wg.Wait()          // дождаться завершения основного воркера

	if err := x.series.Close(); err != nil {
		fmt.Println("close series error:", err)
	}
	if err := panalib.SaveConfig(x.config); err != nil {
		fmt.Println("save config error:", err)
	}
	if err := x.comport.Close(); err != nil {
		fmt.Println("close comport error:", err)
	}
}

func mustPipeListener() net.Listener {

	ln, err := winio.ListenPipe(`\\.\pipe\panalib`, nil)
	if err != nil {
		panic(err)
	}
	return ln

}

func fmtErr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
