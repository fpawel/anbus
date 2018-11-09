package anbus

import "github.com/fpawel/goutils/winapp"

const (
	AppName winapp.AnalitpriborAppName = "anbus"
)
const (
	PeerExeName = "anbusui.exe"
)

func DataFileName() string {
	return AppName.DataFileName("series.sqlite")
}
