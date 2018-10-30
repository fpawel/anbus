package work

import "github.com/pkg/errors"

type CmdSvc struct {
	w *worker
}

func (x *CmdSvc) Perform(v [1]string, _ *struct{}) error {
	c, err := parseTxtCmd(v[0])
	if err != nil {
		return errors.Wrap(err, v[0])
	}
	if c.tokens()[0] == "EXIT" {
		if err := x.w.rpcWnd.Close(); err != nil {
			return errors.Wrap(err, "close window: unexpected error")
		}
	}
	if r, err := c.parseModbusRequest(); err == nil {
		x.w.chModbusRequest <- r
		return nil
	} else {
		return errors.Wrap(err, v[0])
	}
}
