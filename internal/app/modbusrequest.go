package app

import (
	"context"
	"github.com/ansel1/merry"
	"github.com/fpawel/anbus/internal/api/notify"
	"github.com/fpawel/anbus/internal/cfg"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gohelp/myfmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type request struct {
	modbus.Request
	all    bool
	source string
}

func (x txtCmd) parseRequest() (request, error) {

	r := request{
		source: x.source,
	}

	switch x.name() {
	case "ALL":
		r.all = true
	default:
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
		var err error
		if r.Request, err = parseW32(r.Addr, xs[1:]); err != nil {
			return r, errors.Wrap(err, "запись в регистр 32")
		}
		return r, nil
	default:
		if cmdCode, err := parseHexByte(xs[0]); err == nil {
			r.ProtoCmd = modbus.ProtoCmd(cmdCode)
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

func parseW32(addr modbus.Addr, xs []string) (modbus.Request, error) {

	cmd, err := strconv.Atoi(xs[0])
	if err != nil {
		return modbus.Request{}, merry.Append(err, "код кодманды")
	}
	if cmd < 0 || cmd > 0xFFFF {
		return modbus.Request{}, merry.New("код кодманды должен быть от 0 до 0xFFFF")
	}

	v, err := strconv.ParseFloat(strings.Replace(xs[1], ",", ".", -1), 64)
	if err != nil {
		return modbus.Request{}, merry.Append(err, "значение аргумента")
	}
	return modbus.NewWrite32BCDRequest(addr, 0x010, modbus.DevCmd(cmd), v), nil
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

func (r request) perform(ctxWork context.Context) {
	if r.Addr == 0 || !r.all {
		r.performAddr(ctxWork)
		return
	}
	for _, p := range cfg.Get().Nodes() {
		r := r
		r.Addr = p.Addr
		r.performAddr(ctxWork)
	}
	return
}

func (r request) performAddr(ctx context.Context) {
	t := time.Now()

	if r.Addr == 0 {
		if _, err := comPort.Write(log, ctx, r.Bytes()); err != nil {
			notify.WriteConsoleErrorf(nil, "% X : %v", r.Bytes(), err.Error())
		} else {
			notify.WriteConsoleInfof(nil, "% X", r.Bytes())
		}
		return
	}
	response, err := comPort.GetResponse(log, ctx, r.Bytes(), func(_ []byte, _ []byte) (string, error) {
		return "", nil
	})
	strTime := myfmt.FormatDuration(time.Since(t))
	if err == nil {
		notify.WriteConsoleInfof(nil, "% X -> % X %s", r.Bytes(), response, strTime)
	} else {
		notify.WriteConsoleErrorf(nil, "% X : %v %s", r.Bytes(), err, strTime)
	}
}
