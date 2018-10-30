package work

import (
	"context"
	"fmt"
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/pkg/errors"
	"time"
)

func (x *worker) main() {
	var va anbus.VarAddr
	for !x.flagClose {
		sets := x.sets.Config()
		select {
		case r := <-x.chModbusRequest:
			x.notifyConsoleInfo(r.source)
			if x.prepareComport(sets) {
				x.getResponse(r, sets)
				continue
			}
		default:
			va = sets.NextVarAddr(va)
			if va.Place >= 0 && x.prepareComport(sets) {
				if v, ok := x.doReadVar(va); ok && sets.SaveMin > 0 {
					x.series.AddRecord(va.Addr, va.Var, v)
				}
				continue
			}
		}
		time.Sleep(time.Second)
	}
}

func (x *worker) prepareComport(sets anbus.Config) bool {

	comportConfig := x.comport.Config()

	if comportConfig.Uart != sets.Comport.Uart {
		x.comport.SetUartConfig(sets.Comport.Uart)
	}
	if comportConfig.Serial != sets.Comport.Serial {
		if err := x.comport.Close(); err != nil {
			fmt.Println("close COMPORT error:", err)
		}
	}

	if !x.comport.Opened() {
		if err := x.comport.OpenWithConfig(sets.Comport); err != nil {
			x.notifyStatusError("%v", err)
			return false
		}
	}
	return true
}

func (x *worker) getResponse(r modbusRequest, cfg anbus.Config) {

	doAddr := func() {
		if _, err := x.comport.GetResponse(r.Bytes()); err == nil {
			x.notifyConsoleInfo(x.comport.Dump())
		} else {
			x.notifyConsoleError("%s: %v", x.comport.Dump(), err)
		}
	}

	if r.all {
		for _, p := range cfg.Places {
			if !p.Unchecked {
				r.Addr = p.Addr
				doAddr()
			}
		}
		return
	}

	if r.Addr > 0 {
		doAddr()
		return
	}

	if _, err := x.comport.Write(r.Bytes()); err != nil {
		x.notifyConsoleInfo(err.Error())
	} else {
		x.notifyConsoleInfo("% X : BROADCAST", r.Bytes())
	}

}

func (x *worker) doReadVar(va anbus.VarAddr) (float64, bool) {

	value, err := modbus.Read3BCD(x.comport, va.Addr, va.Var)
	if err == context.DeadlineExceeded {
		err = errors.New("нет ответа")
	}
	if err != nil {
		err = errors.New(err.Error() + ": " + x.comport.Dump())
	}

	x.rpcWnd.NotifyParam(msgReadVar, struct {
		Place, VarIndex int
		Value           float64
		Error           string
	}{
		va.Place, va.VarIndex, value, fmtErr(err),
	})
	return value, err == nil
}
