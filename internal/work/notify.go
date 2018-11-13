package work

import "fmt"

func (x *worker) notifyConsoleInfo(format string, a ...interface{}) {
	x.notify(msgInfo, msgConsole, format, a...)
}

func (x *worker) notifyConsoleError(format string, a ...interface{}) {
	x.notify(msgErr, msgConsole, format, a...)
}

func (x *worker) notifyStatusInfo(format string, a ...interface{}) {
	x.notify(msgInfo, msgStatus, format, a...)
}

func (x *worker) notifyStatusError(format string, a ...interface{}) {
	x.notify(msgErr, msgStatus, format, a...)
}

func (x *worker) notify(level MsgTextLevel, kind MsgTextKind, format string, a ...interface{}) {
	x.notifyWindow.NotifyParam(msgText, struct {
		Level MsgTextLevel
		Kind  MsgTextKind
		Text  string
	}{
		level,
		kind,
		fmt.Sprintf(format, a...),
	})
}
