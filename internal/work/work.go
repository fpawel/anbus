package work

import (
	"errors"
	"fmt"
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/fpawel/anbus/internal/notify"
	"github.com/fpawel/anbus/internal/panalib"
	"time"
)

type modbusRequest struct {
	modbus.Request
	all    bool
	source string
}

func (x *worker) mainWork() {
	var va panalib.VarAddr
	for !x.flagClose {
		cfg := x.safeGetConfig()
		select {
		case r := <-x.chModbusRequest:
			x.window.SendConsoleInfo(r.source)
			if x.prepareComport(cfg) {
				x.doModbus(r, cfg)
				continue
			}
		default:
			va = cfg.NextVarAddr(va)
			if va.Place >= 0 && x.prepareComport(cfg) {
				if v, ok := x.doReadVar(va); ok {
					x.series.AddRecord(va.Addr, va.Var, v, time.Minute*time.Duration(cfg.SaveMin))
				}
				continue
			}
		}
		time.Sleep(time.Second)
	}
}

func (x *worker) prepareComport(cfg panalib.Config) bool {

	comportConfig := x.comport.Config()

	if comportConfig.Uart != cfg.Comport.Uart {
		x.comport.SetUartConfig(cfg.Comport.Uart)
	}
	if comportConfig.Serial != cfg.Comport.Serial {
		if err := x.comport.Close(); err != nil {
			fmt.Println("close COMPORT error:", err)
		}
	}

	if !x.comport.Opened() {
		if err := x.comport.OpenWithConfig(cfg.Comport); err != nil {
			x.window.SendStatusError("%v", err)
			return false
		}
	}
	return true

}

func (x *worker) doModbusAddr(r modbusRequest) {
	if _, err := x.comport.GetResponse(r.Bytes()); err == nil {
		x.window.SendConsoleInfo(x.comport.Dump())
	} else {
		x.window.SendConsoleError(x.comport.Dump())
	}
}

func (x *worker) doModbus(r modbusRequest, cfg panalib.Config) {

	if r.all {
		for _, p := range cfg.Places {
			if !p.Unchecked {
				r.Addr = p.Addr
				x.doModbusAddr(r)
			}
		}
		return
	}

	if r.Addr > 0 {
		x.doModbusAddr(r)
		return
	}

	if _, err := x.comport.Write(r.Bytes()); err != nil {
		x.window.SendConsoleInfo(err.Error())
	} else {
		x.window.SendConsoleInfo("< % X : широковещательное сообщение", r.Bytes())
	}

}

func (x *worker) doReadVar(va panalib.VarAddr) (float64, bool) {

	value, err := modbus.Read3BCD(x.comport, va.Addr, va.Var)
	if err != nil {
		err = errors.New(x.comport.Dump())
	}

	r := struct {
		Place, VarIndex int
		Value           float64
		Error           string
	}{
		va.Place, va.VarIndex, value, fmtErr(err),
	}
	x.window.SendMsgJSON(notify.MsgReadVar, &r)
	return value, err == nil
}
