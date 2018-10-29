package work

type Msg uintptr

const (
	msgStatusInfo uintptr = iota
	msgStatusError
	msgConsoleInfo
	msgConsoleError
	msgReadVar
)

type TextMsgLevel string

const (
	msgInfo TextMsgLevel = "info"
	msgErr  TextMsgLevel = "error"
)
