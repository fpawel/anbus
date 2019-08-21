package app

import (
	"context"
	"github.com/ansel1/merry"
	"github.com/fpawel/anbus/internal/api/notify"
	"github.com/fpawel/anbus/internal/api/types"
	"github.com/fpawel/anbus/internal/cfg"
	"github.com/fpawel/comm"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/dseries"
	"time"
)

func work() {
	var n cfg.Node //сетевой объект опроса
	for {
		select {

		case <-chConfigChanged:
			log.ErrIfFail(comPort.Close)

		case <-ctxApp.Done():
			return // работа приложения прервана, выход

		case task := <-chTasks:
			task() // выполнить дополнителную задачу

		default:
			config := cfg.Get()
			// выполние основной работы
			// вычислить следующий сетевой объект опроса
			if n = config.NextNode(n); n.Place < 0 {
				// не заданы сетевые объекты опроса
				pause(time.Second)
				continue
			}
			r := types.ReadVar{
				Place:    n.Place,
				VarIndex: n.VarIndex,
				VarCode:  n.VarCode,
				Addr:     n.Addr,
				VarName:  config.VarNameByCode(n.VarCode),
			}
			var err error

			r.Value, err = modbus.Read3BCD(log, ctxApp, comPort, n.Addr, n.VarCode)
			if err == nil {
				// считано новое значение, отправить оповещение о нём
				notify.ReadVar(nil, r)
				// если предыдущее сохранённое значение было сохранено более 5 минут назад,
				// создать новую пачку графиков
				if time.Since(dseries.UpdatedAt()) > time.Minute*5 {
					dseries.CreateNewBucket("anbus")
					notify.NewSeries(log.Info)
				}
				// сохранить новое значение в базе данных графиков
				dseries.AddPoint(n.Addr, n.VarCode, r.Value)
				continue
			}
			if merry.Is(err, context.Canceled) {
				return // работа приложения прервана, выход
			}
			if isDeviceError(err) {
				// произошла ошибка протокола либо ответ от данного адреса не был получен
				r.Error = err.Error()
				go notify.ReadVar(nil, r)
				continue
			}
			// произошёла ошибка СОМ порта
			notify.WorkError(log.PrintErr, err.Error())
			pause(time.Second)
		}
	}
}

func isDeviceError(err error) bool {
	return merry.Is(err, comm.Err) || merry.Is(err, context.DeadlineExceeded)
}

func pause(d time.Duration) {
	timer := time.NewTimer(d)
	for {
		select {
		case <-timer.C:
			return
		case <-ctxApp.Done():
			timer.Stop()
			return
		}
	}
}
