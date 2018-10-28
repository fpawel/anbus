package notify

type Msg uintptr

const (
	msgStatusText Msg = iota
	msgConsoleText
	MsgReadVar
)

type TextMsgLevel bool

const (
	msgInfo TextMsgLevel = true
	msgErr  TextMsgLevel = false
)
