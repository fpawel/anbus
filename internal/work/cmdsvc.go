package work

import (
	"github.com/pkg/errors"
	"strings"
)

type CmdSvc struct {
	w *worker
}

func (x *CmdSvc) Perform(v [1]string, _ *struct{}) error {
	c, err := parseTxtCmd(v[0])
	if err != nil {
		return errors.Wrap(err, v[0])
	}
	switch strings.ToUpper(c.name()) {
	case "EXIT":
		if !x.w.rpcWnd.CloseWindow() {
			return errors.New("can not close rpc window")
		}
		return nil
	default:
		if r, err := c.parseModbusRequest(); err == nil {
			x.w.chModbusRequest <- r
			return nil
		}
	}
	return errors.Errorf("нет такой команды: %q", c.name())
}
