package types

import "github.com/fpawel/comm/modbus"

type ReadVar struct {
	Place,
	VarIndex int
	Value float64
	Error string
}

type AddrError struct {
	Addr  modbus.Addr
	Place int
	Error string
}

type EmptyRecord struct{}
