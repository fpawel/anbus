package work

import (
	"context"
	"fmt"
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/pkg/errors"
	"time"
)

func (x *worker) work() {
	var va anbus.VarAddr
	for {
		select {
		case <-x.ctx.Done():
			return
		case r := <-x.chModbusRequest:
			cfg := x.sets.Config()
			x.notifyConsoleInfo(r.source)
			if x.prepareComport(cfg) {
				x.getResponse(r, cfg)
				continue
			}
		default:
			cfg := x.sets.Config()
			va = cfg.NextVarAddr(va)
			if va.Place >= 0 && x.prepareComport(cfg) {
				if v, ok := x.doReadVar(va, cfg); ok && cfg.SaveSeries {
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
		if err := x.comport.Open(sets.Comport, x.ctx); err != nil {
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

func (x *worker) doReadVar(va anbus.VarAddr, cfg anbus.Config) (float64, bool) {

	value, err := modbus.Read3BCD(x.comport, va.Addr, va.Var)
	if err == context.DeadlineExceeded {
		err = errors.New("нет ответа")
	}
	if err != nil {
		err = errors.New(err.Error() + ": " + x.comport.Dump())
	}

	x.notifyWindow.NotifyJson(msgReadVar, struct {
		Place, VarIndex int
		Value           float64
		Error           string
	}{
		va.Place, va.VarIndex, value, fmtErr(err),
	})

	if cfg.DumpComport {
		s := time.Now().Format("15:04:05.000")
		if err == nil {
			x.notifyConsoleInfo("%s %s %v", s, x.comport.Dump(), value)
		} else {
			x.notifyConsoleError("%s %v", s, err)
		}
	}
	return value, err == nil
}
