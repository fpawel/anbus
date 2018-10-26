package svc

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/settings"
	"github.com/fpawel/goutils/serial/modbus"
	"time"
)

type SetsSvc struct {
	sets *anbus.Sets
}

func NewSetsSvc(sets *anbus.Sets) *SetsSvc {
	return &SetsSvc{sets}
}

func (x *SetsSvc) UserConfig(_ struct{}, res *settings.Config) error {
	*res = x.sets.UserConfig()
	return nil
}

func (x *SetsSvc) SetSeriesSaveMin(value [1]int, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.SaveMin = value[0]
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) SetMaxAttemptRead(value [1]int, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Comport.Uart.MaxAttemptsRead = value[0]
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) SetReadTimeoutMillis(value [1]int, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Comport.Uart.ReadTimeout = time.Microsecond * time.Duration(value[0])
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) SetReadByteTimeoutMillis(value [1]int, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Comport.Uart.ReadByteTimeout = time.Microsecond * time.Duration(value[0])
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) SetBaud(value [1]int, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Comport.Serial.Baud = value[0]
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) SetPortName(value [1]string, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Comport.Serial.Name = value[0]
	x.sets.SetConfig(cfg)
	return nil
}

type SetVarArg struct {
	Index int
	Var   modbus.Var
}

func (x *SetsSvc) SetVar(v SetVarArg, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Vars[v.Index].Var = v.Var
	x.sets.SetConfig(cfg)
	return nil
}

type SetAddrArg struct {
	Place int
	Addr  modbus.Addr
}

func (x *SetsSvc) SetAddr(v SetAddrArg, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Places[v.Place].Addr = v.Addr
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) Toggle(struct{}, *struct{}) error {
	cfg := x.sets.Config()

	v := len(cfg.NetworkItems()) > 0
	for i := range cfg.Places {
		cfg.Places[i].Unchecked = v
	}
	for i := range cfg.Vars {
		cfg.Vars[i].Unchecked = v
	}
	x.sets.SetConfig(cfg)
	return nil
}

type SetVarCheckedArg struct {
	VarIndex int
	Checked  bool
}

func (x *SetsSvc) SetVarChecked(v SetVarCheckedArg, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Vars[v.VarIndex].Unchecked = !v.Checked
	x.sets.SetConfig(cfg)
	return nil
}

type SetPlaceCheckedArg struct {
	Place   int
	Checked bool
}

func (x *SetsSvc) SetPlaceChecked(v SetPlaceCheckedArg, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Places[v.Place].Unchecked = !v.Checked
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) AddVar(struct{}, *struct{}) error {
	cfg := x.sets.Config()
	cfg.Vars = append(cfg.Vars, anbus.Var{Unchecked: true})
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) DelVar(struct{}, *struct{}) error {

	cfg := x.sets.Config()
	if len(cfg.Vars) > 1 {
		cfg.Vars = cfg.Vars[:len(cfg.Vars)-1]
	}
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) AddPlace(struct{}, *struct{}) error {
	cfg := x.sets.Config()
	cfg.Places = append(cfg.Places, anbus.Place{Unchecked: true})
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) DelPlace(struct{}, *struct{}) error {

	cfg := x.sets.Config()
	if len(cfg.Places) > 1 {
		cfg.Places = cfg.Places[:len(cfg.Places)-1]
	}
	x.sets.SetConfig(cfg)
	return nil
}
