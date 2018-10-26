package data

import (
	"github.com/fpawel/goutils/serial/modbus"
	"time"
)

type Record struct {
	CreatedAt    time.Time
	Var          modbus.Var
	Addr         modbus.Addr
	Value        float64
	CreatedAtStr string
}
