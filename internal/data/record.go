package data

import (
	"github.com/fpawel/goutils/serial/modbus"
	"time"
)

type Record struct {
	StoredAt time.Time
	Var      modbus.Var
	Addr     modbus.Addr
	Value    float64
}
