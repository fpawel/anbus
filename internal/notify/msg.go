package notify

type Msg uintptr

const (
	MsgUserConfig Msg = iota
	MsgNetwork
	msgStatusText
	msgConsoleText
	MsgReadVar
)

type TextMsgLevel bool

const (
	msgInfo TextMsgLevel = true
	msgErr  TextMsgLevel = false
)
