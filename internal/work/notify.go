package work

func (x *worker) notifyConsoleInfo(format string, a ...interface{}) {
	x.rpcWnd.Notifyf(msgConsoleInfo, format, a...)
}

func (x *worker) notifyConsoleError(format string, a ...interface{}) {
	x.rpcWnd.Notifyf(msgConsoleError, format, a...)
}

func (x *worker) notifyStatusInfo(format string, a ...interface{}) {
	x.rpcWnd.Notifyf(msgStatusInfo, format, a...)
}

func (x *worker) notifyStatusError(format string, a ...interface{}) {
	x.rpcWnd.Notifyf(msgStatusError, format, a...)
}
