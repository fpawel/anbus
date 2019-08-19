package app

import (
	"github.com/ansel1/merry"
)

type MainSvc struct{}

func (x *MainSvc) PerformTextCommand(v [1]string, _ *struct{}) error {
	c, err := parseTxtCmd(v[0])
	if err != nil {
		return merry.Append(err, v[0])
	}
	r, err := c.parseModbusRequest()
	if err != nil {
		return err
	}
	chRequest <- r
	return nil
}
