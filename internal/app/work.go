package app

import (
	"context"
	"github.com/fpawel/anbus/internal/cfg"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/hako/durafmt"
	"github.com/pkg/errors"
	"time"
)

func (x *worker) work() {
	var va cfg.VarAddr
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

func (x *worker) prepareComport(sets cfg.Config) bool {

	comportConfig := x.comport.Config()

	if comportConfig.Name != sets.ComportName || comportConfig.Baud != sets.ComportBaud {
		x.comport.Close()
	}

	if !x.comport.Opened() {
		if err := x.comport.Open(sets.ComportName, sets.ComportBaud); err != nil {
			x.notifyStatusError("%v", err)
			return false
		}
	}
	return true
}

func (x *worker) getResponse(r modbusRequest, cfg cfg.Config) {

	doAddr := func() {
		t := time.Now()
		if response, err := x.comport.GetResponse(r.Bytes(), x.sets.Config().Comm, context.Background(),
			func(request []byte, response []byte) error{
				return nil
			}); err == nil {
			x.notifyConsoleInfo( "% X -> % X, %s", r.Bytes(), response,
				durafmt.Parse(time.Since(t)))
		} else {
			x.notifyConsoleInfo( "% X : %v, %s", r.Bytes(), err,
				durafmt.Parse(time.Since(t)))
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

func (x *worker) doReadVar(va cfg.VarAddr, cfg cfg.Config) (float64, bool) {

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
