package work

import (
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type txtCmd struct {
	source string
	name   string
	xs     []string
}

func parseTxtCmd(sourceStr string) (txtCmd, error) {
	xs := strings.Split(sourceStr, " ")
	if len(xs) == 0 {
		return txtCmd{}, errors.New("команда не задана")
	}
	for i := range xs {
		xs[i] = strings.ToUpper(xs[i])
	}
	return txtCmd{
		source: sourceStr,
		name:   xs[0],
		xs:     xs[1:],
	}, nil
}

func (x txtCmd) parseRequest() (request, error) {

	r := request{
		source: x.source,
	}
	if x.name == "ALL" {
		r.all = true
	} else {
		if b, err := parseHexByte(x.name); err == nil {
			r.Addr = modbus.Addr(b)
		} else {
			return r, errors.Wrap(err, "адрес модбас")
		}
	}

	if len(x.xs) == 0 {
		return r, errors.New("не указан код команды модбас")
	}
	if cmdCode, err := parseHexByte(x.xs[0]); err == nil {
		r.ProtocolCommandCode = modbus.ProtocolCommandCode(cmdCode)
	} else {
		return r, errors.Wrap(err, "код команды модбас")
	}
	x.xs = x.xs[1:]
	for i, s := range x.xs {
		if b, err := parseHexByte(s); err == nil {
			r.Data = append(r.Data, byte(b))
		} else {
			return r, errors.Wrapf(err, "байт данных в позиции %d", i)
		}
	}
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
