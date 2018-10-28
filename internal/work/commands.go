package work

import (
	"github.com/fpawel/goutils"
	"github.com/fpawel/goutils/copydata"
	"github.com/lxn/win"
	"unsafe"
)

func (x *worker) onCommand(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {

	switch msg {

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
