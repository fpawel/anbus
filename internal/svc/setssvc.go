package svc

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/settings"
	"github.com/fpawel/goutils/serial/modbus"
	"strconv"
	"time"
)

type SetsSvc struct {
	sets *anbus.Sets
}

func NewSetsSvc(sets *anbus.Sets) *SetsSvc {
	return &SetsSvc{sets}
}

func (x *SetsSvc) Network(_ struct{}, res *anbus.Network) error {
	cfg := x.sets.Config()
	*res = cfg.Network
	return nil
}

func (x *SetsSvc) UserConfig(_ struct{}, res *settings.Config) error {
	*res = x.sets.UserConfig()
	return nil
}

func (x *SetsSvc) SetValue(p settings.PropertyValue, _ *struct{}) error {
	cfg := x.sets.Config()

	switch p.Section {
	case "comport":
		switch p.Name {
		case "name":
			cfg.Comport.Serial.Name = p.Value
		case "baud":
			n, err := strconv.Atoi(p.Value)
			if err != nil {
				return err
			}
			cfg.Comport.Serial.Baud = n
		case "timeout":
			n, err := strconv.Atoi(p.Value)
			if err != nil {
				return err
			}
			cfg.Comport.Uart.ReadTimeout = time.Millisecond * time.Duration(n)
		case "timeout_byte":
			n, err := strconv.Atoi(p.Value)
			if err != nil {
				return err
			}
			cfg.Comport.Uart.ReadByteTimeout = time.Millisecond * time.Duration(n)

		case "max_attempts_read":
			n, err := strconv.Atoi(p.Value)
			if err != nil {
				return err
			}
			cfg.Comport.Uart.MaxAttemptsRead = n
		}
	case "chart":
		switch p.Name {
		case "save_series":
			n, err := strconv.ParseBool(p.Value)
			if err != nil {
				return err
			}
			cfg.SaveSeries = n

		}
	}

	x.sets.SetConfig(cfg)
	return nil
}

type SetVarArg struct {
	Index int        `json:"index"`
	Var   modbus.Var `json:"var"`
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

func (x *SetsSvc) Toggle(_ struct{}, res *anbus.Network) error {
	cfg := x.sets.Config()

	v := len(cfg.NetworkItems()) > 0
	for i := range cfg.Places {
		cfg.Places[i].Unchecked = v
	}
	for i := range cfg.Vars {
		cfg.Vars[i].Unchecked = v
	}
	x.sets.SetConfig(cfg)

	*res = cfg.Network
	return nil
}

func (x *SetsSvc) ToggleVar(v [1]int, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Vars[v[0]].Unchecked = !cfg.Vars[v[0]].Unchecked
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) TogglePlace(v [1]int, _ *struct{}) error {
	cfg := x.sets.Config()
	cfg.Places[v[0]].Unchecked = !cfg.Places[v[0]].Unchecked
	x.sets.SetConfig(cfg)
	return nil
}

func (x *SetsSvc) AddVar(_ struct{}, res *anbus.Network) error {
	cfg := x.sets.Config()
	cfg.Vars = append(cfg.Vars, anbus.Var{Unchecked: true})
	x.sets.SetConfig(cfg)
	*res = cfg.Network
	return nil
}

func (x *SetsSvc) DelVar(_ struct{}, res *anbus.Network) error {

	cfg := x.sets.Config()
	if len(cfg.Vars) > 1 {
		cfg.Vars = cfg.Vars[:len(cfg.Vars)-1]
	}
	x.sets.SetConfig(cfg)
	*res = cfg.Network
	return nil
}

func (x *SetsSvc) AddPlace(_ struct{}, res *anbus.Network) error {
	cfg := x.sets.Config()
	cfg.Places = append(cfg.Places, anbus.Place{Unchecked: true})
	x.sets.SetConfig(cfg)
	*res = cfg.Network
	return nil
}

func (x *SetsSvc) DelPlace(_ struct{}, res *anbus.Network) error {

	cfg := x.sets.Config()
	if len(cfg.Places) > 1 {
		cfg.Places = cfg.Places[:len(cfg.Places)-1]
	}
	x.sets.SetConfig(cfg)
	*res = cfg.Network
	return nil
}
