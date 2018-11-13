package work

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/chart"
	"github.com/fpawel/goutils/copydata"
	"github.com/fpawel/goutils/serial/comport"
	"net"
)

type worker struct {
	notifyWindow    *copydata.NotifyWindow
	comport         *comport.Port
	sets            *anbus.Sets
	series          *chart.Series
	chModbusRequest chan modbusRequest
	ln              net.Listener
	chartSvc        *ChartSvc
}
