package app

import (
	"fmt"
	"github.com/fpawel/anbus/internal/cfg"
	"github.com/fpawel/comm"
	"github.com/fpawel/comm/modbus"
)

type taskRunner struct{}

func (_ taskRunner) Write32(addr int, cmd modbus.DevCmd, value float64) error {
	r := modbus.NewWrite32BCDRequest(0, 0x10, cmd, value)
	if addr == 0 {
		return sendBroadcast(r)
	}
	return comm.WithLogAnswers(func() error {
		for _, addr := range getAddresses(addr) {
			r.Addr = addr
			if _, err := r.GetResponse(log, ctxApp, comPort, func(_, _ []byte) (string, error) {
				return "ok write 32", nil
			}); err != nil || !isDeviceError(err) {
				return err
			}
		}
		return nil
	})
}

func (_ taskRunner) Send(addr int, cmd modbus.ProtoCmd, data []byte) error {

	if addr == 0 {
		return sendBroadcast(modbus.Request{
			Addr:     0,
			ProtoCmd: cmd,
			Data:     data,
		})
	}

	return comm.WithLogAnswers(func() error {
		for _, addr := range getAddresses(addr) {
			_, err := modbus.Request{
				Addr:     addr,
				ProtoCmd: cmd,
				Data:     data,
			}.GetResponse(log, ctxApp, comPort, func(_, _ []byte) (string, error) {
				return "", nil
			})
			if err != nil || !isDeviceError(err) {
				return err
			}
		}
		return nil
	})
}

func sendBroadcast(r modbus.Request) error {
	if _, err := comPort.Write(log, ctxApp, r.Bytes()); err != nil {
		return err
	}
	log.Info(fmt.Sprintf("`% X`", r.Bytes()))
	return nil
}

func getAddresses(addr int) []modbus.Addr {
	var addresses []modbus.Addr
	if addr > 0 {
		return append(addresses, modbus.Addr(addr))
	}
	for _, p := range cfg.Get().Places {
		if p.Check {
			addresses = append(addresses, p.Addr)
		}
	}
	return addresses
}
