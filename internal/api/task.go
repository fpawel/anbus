package api

import (
	"github.com/ansel1/merry"
	"github.com/fpawel/comm/modbus"
	"strconv"
	"strings"
)

type TaskRunner interface {
	Send(int, modbus.ProtoCmd, []byte) error
	Write32(int, modbus.DevCmd, float64) error
}

type TaskSvc struct {
	R TaskRunner
}

func (x TaskSvc) Write32(r struct {
	Addr  int
	Cmd   modbus.DevCmd
	Value float64
}, _ *struct{}) error {
	return x.R.Write32(r.Addr, r.Cmd, r.Value)
}

func (x TaskSvc) Send(r struct {
	Addr  int
	Cmd   modbus.ProtoCmd
	Bytes string
}, _ *struct{}) error {
	bs, err := parseHexBytes(r.Bytes)
	if err != nil {
		return err
	}
	return x.R.Send(r.Addr, r.Cmd, bs)
}

func parseHexBytes(s string) ([]byte, error) {
	var xs []byte
	s = strings.TrimSpace(s)
	for i, strB := range strings.Split(s, " ") {
		v, err := strconv.ParseUint(strB, 16, 8)
		if err != nil {
			return nil, merry.Appendf(err, "поз.%d: %q", i+1, strB)
		}
		if v < 0 || v > 0xff {
			return nil, merry.Errorf("поз.%d: %q: ожижалось шестнадцатиричное число от 0 до FF", i+1, strB)
		}
		xs = append(xs, byte(v))
	}
	return xs, nil
}
