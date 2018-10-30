package work

import (
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type modbusRequest struct {
	modbus.Request
	all    bool
	source string
}

func (x txtCmd) parseModbusRequest() (modbusRequest, error) {

	r := modbusRequest{
		source: x.source,
	}

	if x.name() == "ALL" {
		r.all = true
	} else {
		if b, err := parseHexByte(x.tokens()[0]); err == nil {
			r.Addr = modbus.Addr(b)
		} else {
			return r, errors.Wrap(err, "адрес модбас")
		}
	}
	xs := x.tokens()[1:]

	if len(xs) == 0 {
		return r, errors.New("не указан код команды протокола MODBUS RTU")
	}

	switch xs[0] {
	case "W32":
		r.ProtocolCommandCode = 0x10
		xs = xs[1:]
		b, err := parseW32(xs)
		if err != nil {
			return r, errors.Wrap(err, "кодманда регистра 32")
		}
		r.Data = append(r.Data, b...)
	default:
		if cmdCode, err := parseHexByte(xs[0]); err == nil {
			r.ProtocolCommandCode = modbus.ProtocolCommandCode(cmdCode)
		} else {
			return r, errors.Wrap(err, "код команды протокола MODBUS RTU")
		}
		xs = xs[1:]
		for i, s := range xs {
			if b, err := parseHexByte(s); err == nil {
				r.Data = append(r.Data, byte(b))
			} else {
				return r, errors.Wrapf(err, "байт данных в позиции %d", i)
			}
		}
	}

	return r, nil
}

func parseW32(xs []string) ([]byte, error) {
	var r []byte
	if b, err := parseHexByte(xs[0]); err == nil {
		r = append(r, byte(b))
	} else {
		return r, errors.Wrap(err, "старший байт кода кодманды")
	}

	if b, err := parseHexByte(xs[1]); err == nil {
		r = append(r, byte(b))
	} else {
		return r, errors.Wrap(err, "младший байт кода кодманды")
	}

	v, err := strconv.ParseFloat(strings.Replace(xs[2], ",", ".", -1), 64)
	if err != nil {
		return r, err
	}
	r = append(r, modbus.BCD6(v)...)
	return r, nil
}

func parseHexByte(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0, err
	}
	if v < 0 || v > 0xff {
		return 0, errors.Errorf("%q: ожидалось 8 битное число без знака", s)
	}
	return v, nil
}
