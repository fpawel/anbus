package work

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/data"
	"github.com/fpawel/goutils/copydata"
	"github.com/fpawel/goutils/serial/comport"
	"net"
)

type worker struct {
	rpcWnd          *copydata.RPCWindow
	comport         *comport.Port
	sets            *anbus.Sets
	flagClose       bool
	series          *data.Series
	chModbusRequest chan request
	ln              net.Listener
}
