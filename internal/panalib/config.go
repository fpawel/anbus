package panalib

import (
	"encoding/json"
	"github.com/fpawel/goutils/serial/comport"
	"github.com/fpawel/goutils/serial/modbus"
	"io/ioutil"
)

type Config struct {
	Places  []Place        `json:"places"`
	Vars    []Var          `json:"vars"`
	Comport comport.Config `json:"comport"`
	SaveMin int            `json:"save_min"`
}

type Place struct {
	Addr      modbus.Addr `json:"addr"`
	Unchecked bool        `json:",omitempty"`
}

type Var struct {
	Var       modbus.Var `json:"var"`
	Unchecked bool       `json:"unchecked,omitempty"`
}

type VarAddr struct {
	Var      modbus.Var
	VarIndex int
	Addr     modbus.Addr
	Place    int
}

func configFileName() string {
	return AppName.FileName("config.json")
}

func DefaultConfig() Config {
	return Config{
		Comport: comport.DefaultConfig(),
		Places: []Place{
			{Addr: 1},
		},
		Vars: []Var{{}},
	}
}

func OpenConfig() (Config, error) {
	cfg := DefaultConfig()
	b, err := ioutil.ReadFile(configFileName())
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(b, &cfg)
	return cfg, err
}

func SaveConfig(cfg Config) error {
	b, err := json.MarshalIndent(&cfg, "", "    ")
	if err != nil {
		panic(err)
	}
	return ioutil.WriteFile(configFileName(), b, 0666)
}

func (x Config) Toggle() {
	v := len(x.varAddrItems()) > 0
	for i := range x.Places {
		x.Places[i].Unchecked = v
	}
	for i := range x.Vars {
		x.Vars[i].Unchecked = v
	}
}

func (x Config) varAddrItems() (xs []VarAddr) {
	for place, p := range x.Places {
		for varIndex, v := range x.Vars {
			if !p.Unchecked && !v.Unchecked {
				xs = append(xs, VarAddr{v.Var, varIndex, p.Addr, place})
			}
		}
	}
	return
}

func (x Config) NextVarAddr(va VarAddr) VarAddr {
	xs := x.varAddrItems()
	if len(xs) == 0 {
		return VarAddr{Place: -1}
	}
	for i, vb := range xs {
		if vb == va && i < len(xs)-1 {
			return xs[i+1]
		}
	}
	return xs[0]
}
