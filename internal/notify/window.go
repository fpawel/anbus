package notify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fpawel/goutils"
	"github.com/fpawel/goutils/copydata"
	"github.com/fpawel/goutils/winapp"
	"github.com/lxn/win"
)

type Window struct {
	hWnd     win.HWND
	hWndPeer win.HWND
}

func NewWindow(windowProcedure winapp.WindowProcedure) *Window {
	w := &Window{
		hWnd: winapp.NewWindowWithClassName("AnbusServerAppWindow", windowProcedure),
	}
	w.FindPeerWindow()
	return w
}

func (x *Window) Close() error {
	if win.SendMessage(x.hWnd, win.WM_CLOSE, 0, 0) != 0 {
		return errors.New("can not close window")
	}
	return nil
}

func (x *Window) CheckPeerWindow() bool {
	return winapp.IsWindow(x.hWndPeer)
}

func (x *Window) FindPeerWindow() {
	ptrClassName := goutils.MustUTF16PtrFromString("TAnbusMainForm")
	x.hWndPeer = win.FindWindow(ptrClassName, nil)
}

func (x *Window) SendMsg(msg Msg, b []byte) {
	if !winapp.IsWindow(x.hWndPeer) {
		x.FindPeerWindow()
	}
	if winapp.IsWindow(x.hWndPeer) && copydata.SendMessage(x.hWnd, x.hWndPeer, uintptr(msg), b) == 0 {
		x.hWndPeer = 0
	}
}

func (x *Window) SendMsgStr(msg Msg, s string) {
	x.SendMsg(msg, goutils.UTF16FromString(s))
}

func (x *Window) SendMsgJSON(msg Msg, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	x.SendMsgStr(msg, string(b))
}

func (x *Window) SendConsoleInfo(format string, a ...interface{}) {
	x.sendMsgText(msgConsoleText, msgInfo, format, a...)
}

func (x *Window) SendConsoleError(format string, a ...interface{}) {
	x.sendMsgText(msgConsoleText, msgErr, format, a...)
}

func (x *Window) SendStatusInfo(format string, a ...interface{}) {
	x.sendMsgText(msgStatusText, msgInfo, format, a...)
}

func (x *Window) SendStatusError(format string, a ...interface{}) {
	x.sendMsgText(msgStatusText, msgErr, format, a...)
}

func (x *Window) sendMsgText(msg Msg, l TextMsgLevel, format string, a ...interface{}) {
	switch msg {
	case msgStatusText, msgConsoleText:
		x.SendMsgJSON(msg, struct {
			Text string
			Ok   TextMsgLevel
		}{
			fmt.Sprintf(format, a...), l,
		})
	default:
		panic(fmt.Sprintf("%v: must be msgStatusText or msgConsoleText", msg))
	}
}
