package work

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Microsoft/go-winio"
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/anbus/svccfg"
	"github.com/fpawel/anbus/internal/chart"
	"github.com/fpawel/goutils/copydata"
	"github.com/fpawel/goutils/serial/comport"
	"github.com/fpawel/goutils/winapp"
	"github.com/hashicorp/go-multierror"
	"github.com/lxn/win"
	"github.com/pkg/errors"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
)

const (
	pipeName                      = `\\.\pipe\anbus`
	anbusServerAppWindowClassName = "AnbusServerAppWindow"
	peerWindowClassName           = "TAnbusMainForm"
)

func Main(mustRunPeer bool) {

	series := chart.NewSeries()
	x := &worker{
		sets:            openConfig(),
		chModbusRequest: make(chan modbusRequest, 10),
		ln:              mustPipeListener(),
		notifyWindow:    copydata.NewNotifyWindow(anbusServerAppWindowClassName, peerWindowClassName),
		series:          series,
		chartSvc:        &ChartSvc{series},
		comport:         new(comport.Port),
	}

	if mustRunPeer && !winapp.IsWindow(FindPeer()) {
		if err := runPeer(); err != nil {
			panic(err)
		}
	}

	var cancel func()
	x.ctx, cancel = context.WithCancel(context.Background())

	rpcMustRegister(
		svccfg.NewSetsSvc(x.sets),
		&MainSvc{x},
		x.chartSvc)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// цикл rpc
	go func() {
		defer wg.Done()
		defer x.notifyWindow.CloseWindow()
		count := int32(0)
		for {
			switch conn, err := x.ln.Accept(); err {
			case nil:
				go func() {
					atomic.AddInt32(&count, 1)
					jsonrpc2.ServeConnContext(x.ctx, conn)
					if atomic.AddInt32(&count, -1) == 0 && mustRunPeer {
						return
					}
				}()
			case winio.ErrPipeListenerClosed:
				return
			default:
				fmt.Println("rpc pipe error:", err)
				return
			}
		}
	}()

	// цикл компорта
	go func() {
		x.work()
		wg.Done()
	}()

	// цикл оконных сообщений
	runWindowMessageLoop()
	cancel()
	if err := x.ln.Close(); err != nil {
		fmt.Println("close pipe listener error:", err)
	}
	wg.Wait()
	if err := x.Close(); err != nil {
		fmt.Println(err)
	}
}

func (x *worker) Close() (result error) {

	if err := x.series.Close(); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "close DB series"))
	}
	if err := x.sets.Save(); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "save config"))
	}
	if err := x.comport.Close(); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "close comport"))
	}
	for hWnd := FindPeer(); winapp.IsWindow(hWnd); hWnd = FindPeer() {
		if win.SendMessage(hWnd, win.WM_CLOSE, 0, 0) != 0 {
			result = multierror.Append(result, errors.New("can not close peer window"))
		}
	}
	return
}

func runPeer() error {
	const (
		peerAppExe = "anbusui.exe"
	)
	dir := filepath.Dir(os.Args[0])

	if _, err := os.Stat(filepath.Join(dir, peerAppExe)); os.IsNotExist(err) {
		dir = anbus.AppName.Dir()
	}

	cmd := exec.Command(filepath.Join(dir, peerAppExe))
	cmd.Stdout = os.Stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	return cmd.Start()
}

func mustPipeListener() net.Listener {

	ln, err := winio.ListenPipe(pipeName, nil)
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

func rpcMustRegister(rcvrs ...interface{}) {
	for _, rcvr := range rcvrs {
		if err := rpc.Register(rcvr); err != nil {
			panic(err)
		}
	}
}

func FindPeer() win.HWND {
	return winapp.FindWindow(peerWindowClassName)
}

func runWindowMessageLoop() {
	for {
		var msg win.MSG
		if win.GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}
