package settings

type ValueType int

const (
	VtInt ValueType = iota
	VtFloat
	VtString
	VtComportName
	VtBaud
	VtBool
)
