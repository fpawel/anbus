package work

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/data"
	"github.com/fpawel/anbus/internal/notify"
	"github.com/fpawel/goutils/serial/comport"
	"net"
)

type worker struct {
	window          *notify.Window
	comport         *comport.Port
	sets            *anbus.Sets
	flagClose       bool
	series          *data.Series
	chModbusRequest chan request
	ln              net.Listener
}
