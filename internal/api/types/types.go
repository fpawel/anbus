package types

import "github.com/fpawel/comm/modbus"

type ReadVar struct {
	Place    int
	VarIndex int
	VarCode  modbus.Var
	Addr     modbus.Addr
	VarName  string
	Value    float64
	Error    string
}

type EmptyRecord struct{}
