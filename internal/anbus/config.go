package anbus

import (
	"github.com/fpawel/goutils/serial/comport"
)

type Config struct {
	Network
	Comport    comport.Config `json:"comport"`
	SaveSeries bool           `json:"save_series"`
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
