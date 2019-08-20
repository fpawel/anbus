package api

import (
	"github.com/fpawel/anbus/internal/cfg"
	"github.com/fpawel/comm/modbus"
	"github.com/pelletier/go-toml"
)

type ConfigSvc struct{}

func (_ *ConfigSvc) Network(_ struct{}, r *cfg.ConfigNetwork) error {
	*r = cfg.Get().ConfigNetwork
	return nil
}

func (_ *ConfigSvc) SetVar(x struct {
	VarIndex int
	VarCode  modbus.Var
}, _ *struct{}) error {
	c := cfg.Get()
	c.Vars[x.VarIndex].Code = x.VarCode
	cfg.Set(c)
	return nil
}

func (_ *ConfigSvc) SetAddr(x struct {
	Place int
	Addr  modbus.Addr
}, _ *struct{}) error {
	c := cfg.Get()
	c.Places[x.Place].Addr = x.Addr
	cfg.Set(c)
	return nil
}

func (_ *ConfigSvc) ToggleNetwork(_ struct{}, r *cfg.ConfigNetwork) error {
	c := cfg.Get()
	c.ToggleChecked()
	cfg.Set(c)
	*r = c.ConfigNetwork
	return nil
}

func (_ *ConfigSvc) ToggleVar(v [1]int, _ *struct{}) error {
	c := cfg.Get()
	c.Vars[v[0]].Check = !c.Vars[v[0]].Check
	cfg.Set(c)
	return nil
}

func (_ *ConfigSvc) TogglePlace(v [1]int, _ *struct{}) error {
	c := cfg.Get()
	c.Places[v[0]].Check = !c.Places[v[0]].Check
	cfg.Set(c)
	return nil
}

func (_ *ConfigSvc) AddVar(_ struct{}, r *cfg.ConfigNetwork) error {
	c := cfg.Get()
	c.Vars = append(c.Vars, cfg.DevVar{})
	cfg.Set(c)
	*r = c.ConfigNetwork
	return nil
}

func (x *ConfigSvc) DelVar(_ struct{}, r *cfg.ConfigNetwork) error {
	c := cfg.Get()
	if len(c.Vars) > 1 {
		c.Vars = c.Vars[:len(c.Vars)-1]
	}
	cfg.Set(c)
	*r = c.ConfigNetwork
	return nil
}

func (_ *ConfigSvc) AddPlace(_ struct{}, r *cfg.ConfigNetwork) error {
	c := cfg.Get()
	c.Places = append(c.Places, cfg.Place{})
	cfg.Set(c)
	*r = c.ConfigNetwork
	return nil
}

func (_ *ConfigSvc) DelPlace(_ struct{}, r *cfg.ConfigNetwork) error {

	c := cfg.Get()
	if len(c.Places) > 1 {
		c.Places = c.Places[:len(c.Places)-1]
	}
	cfg.Set(c)
	*r = c.ConfigNetwork
	return nil
}

func (_ *ConfigSvc) GetEditConfig(_ struct{}, r *string) error {
	b, err := toml.Marshal(cfg.Get().ConfigEditable)
	if err != nil {
		return err
	}
	*r = string(b)
	return nil
}

func (x *ConfigSvc) SetEditConfig(s [1]string, r *string) error {
	c := cfg.Get()
	if err := toml.Unmarshal([]byte(s[0]), &c.ConfigEditable); err != nil {
		return err
	}
	b, err := toml.Marshal(&c.ConfigEditable)
	if err != nil {
		return err
	}
	*r = string(b)
	cfg.Set(c)
	return nil
}

func (_ *ConfigSvc) Vars(_ struct{}, vars *[]cfg.Var) error {
	*vars = cfg.Get().VarsList
	if len(*vars) == 0 {
		*vars = []cfg.Var{}
	}
	return nil
}
