package work

import (
	"github.com/fpawel/goutils"
	"github.com/fpawel/goutils/copydata"
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/fpawel/anbus/internal/panalib"
	"github.com/lxn/win"
	"unsafe"
)

const (
	cmdPeer = iota + win.WM_USER
	cmdSetVarChecked
	cmdSetPlaceChecked
	cmdAddDelVar
	cmdAddDelPlace
	cmdSetAddr
	cmdSetVar
	cmdToggle
)

func (x *worker) onCommand(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {

	switch msg {

	case cmdPeer:
		x.initPeer()
		return 1

	case cmdSetVarChecked:
		x.muConfig.Lock()
		defer x.muConfig.Unlock()
		x.config.Vars[int(wParam)].Unchecked = lParam != 0
		return 1

	case cmdSetPlaceChecked:
		x.muConfig.Lock()
		defer x.muConfig.Unlock()
		x.config.Places[int(wParam)].Unchecked = lParam != 0
		return 1

	case cmdAddDelVar:
		x.muConfig.Lock()
		if int(wParam) == 0 {
			x.config.Vars = append(x.config.Vars, panalib.Var{Unchecked: true})
		} else {
			if len(x.config.Vars) < 2 {
				x.muConfig.Unlock()
				return 1
			}
			x.config.Vars = x.config.Vars[:len(x.config.Vars)-1]
		}
		x.muConfig.Unlock()
		x.initPeer()
		return 1

	case cmdAddDelPlace:
		x.muConfig.Lock()
		if int(wParam) == 0 {
			x.config.Places = append(x.config.Places, panalib.Place{Unchecked: true})
		} else {
			if len(x.config.Places) < 2 {
				x.muConfig.Unlock()
				return 1
			}
			x.config.Places = x.config.Places[:len(x.config.Places)-1]
		}
		x.muConfig.Unlock()
		x.initPeer()
		return 1

	case cmdSetAddr:
		x.muConfig.Lock()
		defer x.muConfig.Unlock()
		x.config.Places[int(wParam)].Addr = modbus.Addr(lParam)
		return 1

	case cmdSetVar:
		x.muConfig.Lock()
		defer x.muConfig.Unlock()
		x.config.Vars[int(wParam)].Var = modbus.Var(lParam)
		return 1

	case cmdToggle:
		x.muConfig.Lock()
		x.config.Toggle()
		x.muConfig.Unlock()
		x.initPeer()
		return 1

	case win.WM_COPYDATA:
		cdMSG, b16 := copydata.GetData(unsafe.Pointer(lParam))
		data, err := goutils.UTF8FromUTF16(b16)
		if err != nil {
			panic(err)
		}
		x.onCopyData(CommandDataPeer(cdMSG), data)
		return 1

	default:
		return win.DefWindowProc(hWnd, msg, wParam, lParam)
	}
}
