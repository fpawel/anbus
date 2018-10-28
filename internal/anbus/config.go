package anbus

import (
	"github.com/fpawel/goutils/serial/comport"
)

type Config struct {
	Network
	Comport comport.Config `json:"comport"`
	SaveMin int            `json:"save_min"`
}

func configFileName() string {
	return AppName.FileName("config.json")
}

func defaultConfig() Config {
	return Config{
		Comport: comport.DefaultConfig(),
		Network: Network{
			Places: []Place{
				{Addr: 1},
			},
			Vars: []Var{{}},
		},
	}
}