package anbus

import "github.com/fpawel/goutils/serial/modbus"

type Network struct {
	Places []Place `json:"places"`
	Vars   []Var   `json:"vars"`
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

func (x *Network) ToggleChecked() {

	v := len(x.NetworkItems()) > 0
	for i := range x.Places {
		x.Places[i].Unchecked = v
	}
	for i := range x.Vars {
		x.Vars[i].Unchecked = v
	}
}

func (x *Network) NextVarAddr(va VarAddr) VarAddr {
	xs := x.NetworkItems()
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

func (x *Network) NetworkItems() (xs []VarAddr) {
	for place, p := range x.Places {
		for varIndex, v := range x.Vars {
			if !p.Unchecked && !v.Unchecked {
				xs = append(xs, VarAddr{v.Var, varIndex, p.Addr, place})
			}
		}
	}
	return
}
