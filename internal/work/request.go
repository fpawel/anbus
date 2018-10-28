package work

import "github.com/fpawel/goutils/serial/modbus"

type request struct {
	modbus.Request
	all    bool
	source string
}
