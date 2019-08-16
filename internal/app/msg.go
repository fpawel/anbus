package app

type Msg uintptr

const (
	msgText uintptr = iota
	msgReadVar
)

type MsgTextLevel string

type MsgTextKind string

const (
	msgInfo    MsgTextLevel = "info"
	msgErr     MsgTextLevel = "error"
	msgConsole MsgTextKind  = "console"
	msgStatus  MsgTextKind  = "status"
)
