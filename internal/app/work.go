package app

import (
	"context"
	"github.com/ansel1/merry"
	"github.com/fpawel/anbus/internal/api/notify"
	"github.com/fpawel/anbus/internal/cfg"
	"github.com/fpawel/comm/comport"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gohelp/myfmt"
	"github.com/hako/durafmt"
	"github.com/pkg/errors"
	"sync"
	"time"
)

func perform() {
	cancelWorkFunc()
	wgWork.Wait()
	wgWork = sync.WaitGroup{}
	var ctxWork context.Context
	ctxWork, cancelWorkFunc = context.WithCancel(ctxApp)
	wgWork.Add(1)

	go func() {
		defer func() {
			log.ErrIfFail(comPort.Close)
			wgWork.Done()
		}()



		notify.WorkStarted(log)

		err := work(worker)
		if err == nil {
			worker.log.Info("выполнено успешно")
			notify.WorkComplete(worker.log, api.WorkResult{workName, wrOk, "успешно"})
			return
		}

		kvs := merryKeysValues(err)
		if merry.Is(err, context.Canceled) {
			worker.log.Warn("выполнение прервано", kvs...)
			notify.WorkComplete(worker.log, api.WorkResult{workName, wrCanceled, "перервано"})
			return
		}
		worker.log.PrintErr(err, append(kvs, "stack", myfmt.FormatMerryStacktrace(err))...)
		notify.WorkComplete(worker.log, api.WorkResult{workName, wrError, err.Error()})
	}()
}

func work(ctxWork context.Context) error {
	if err := comPort.Open(log, ctxApp); err != nil {
		return err
	}
	var va cfg.Node
	for {
		select {
		case <-ctxApp.Done():
			return nil
		case r := <-chRequest:

			notify.WriteConsole(log, r.source)

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

func prepareComport() bool {

	c :=cfg.Get()

	if c.ComportName != comPort. || comportConfig.Baud != sets.ComportBaud {
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



func doReadVar(n cfg.Node, ctxWork context.Context) (float64, bool) {

	value, err := modbus.Read3BCD(log, ctxWork, comPort, n.Addr, n.VarCode)
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
