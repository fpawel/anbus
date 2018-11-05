package work

import (
	"bytes"
	"fmt"
	"github.com/Microsoft/go-winio"
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/data"
	"github.com/fpawel/anbus/internal/svc"
	"github.com/fpawel/goutils/copydata"
	"github.com/fpawel/goutils/serial/comport"
	"github.com/fpawel/goutils/winapp"
	"github.com/lxn/win"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

const (
	pipeName                      = `\\.\pipe\anbus`
	anbusServerAppWindowClassName = "AnbusServerAppWindow"
	peerWindowClassName           = "TAnbusMainForm"
)

func Main(mustRunPeer bool) {

	x := &worker{
		sets:            openConfig(),
		chModbusRequest: make(chan modbusRequest, 10),
		ln:              mustPipeListener(),
		rpcWnd:          copydata.NewRPCWindow(anbusServerAppWindowClassName, peerWindowClassName),
	}

	if mustRunPeer && !winapp.IsWindow(findPeer()) {
		if err := runPeer(); err != nil {
			panic(err)
		}
	}

	x.comport = comport.NewPortWithConfig(x.sets.Config().Comport)
	x.series = data.NewSeries()

	rpcMustRegister(
		svc.NewSetsSvc(x.sets),
		&CmdSvc{x},
		x.series.Buckets())

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
		x.main()
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

	for hWnd := findPeer(); winapp.IsWindow(hWnd); hWnd = findPeer() {
		win.SendMessage(hWnd, win.WM_CLOSE, 0, 0)
	}
}

func runPeer() error {
	const (
		peerAppExe = "anbusui.exe"
	)
	dir := filepath.Dir(os.Args[0])

	if _, err := os.Stat(filepath.Join(dir, peerAppExe)); os.IsNotExist(err) {
		dir = anbus.AppName.Dir()
	}

	cmd := exec.Command(filepath.Join(dir, peerAppExe), "-must-close-server")
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

func findPeer() win.HWND {
	return winapp.FindWindow(peerWindowClassName)
}
