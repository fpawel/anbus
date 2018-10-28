package work

import (
	"fmt"
	"github.com/Microsoft/go-winio"
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/data"
	"github.com/fpawel/anbus/internal/notify"
	"github.com/fpawel/anbus/internal/svc"
	"github.com/fpawel/goutils/serial/comport"
	"github.com/lxn/win"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"net"
	"net/rpc"
	"sync"
)

const PipeName = `\\.\pipe\anbus`

func Main() {
	x := &worker{
		sets:            openConfig(),
		chModbusRequest: make(chan modbusRequest, 10),
		series:          data.NewSeries(),
		ln:              mustPipeListener(),
	}
	x.comport = comport.NewPortWithConfig(x.sets.Config().Comport)
	x.window = notify.NewWindow(x.onCommand)

	if err := rpc.Register(svc.NewSetsSvc(x.sets)); err != nil {
		panic(err)
	}

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
	if err := x.sets.Save(); err != nil {
		fmt.Println("save sets error:", err)
	}
	if err := x.comport.Close(); err != nil {
		fmt.Println("close comport error:", err)
	}

}

func mustPipeListener() net.Listener {

	ln, err := winio.ListenPipe(PipeName, nil)
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

func openConfig() *anbus.Sets {
	cfg, errCfg := anbus.OpenSets()
	if errCfg != nil {
		fmt.Println("sets:", errCfg)
	}
	return cfg
}

//type debugReadWriteCloser struct {
//	conn net.Conn
//}
//
//func (x *debugReadWriteCloser) Write(p []byte) (int, error) {
//	n,err := x.conn.Write(p)
//	return n,err
//}
//
//func (x *debugReadWriteCloser) Read(p []byte) (int, error) {
//	n,err := x.conn.Read(p)
//	return n,err
//}
//
//func (x *debugReadWriteCloser) Close() error {
//	return x.conn.Close()
//}
