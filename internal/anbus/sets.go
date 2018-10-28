package anbus

import (
	"encoding/json"
	"fmt"
	"github.com/fpawel/anbus/internal/settings"
	"io/ioutil"
	"strconv"
	"sync"
)

type Sets struct {
	cfg Config
	mu  sync.Mutex
}

func OpenSets() (*Sets, error) {
	sets := defaultConfig()
	b, err := ioutil.ReadFile(configFileName())
	if err == nil {
		err = json.Unmarshal(b, &sets)
	}

	return &Sets{cfg: sets}, err
}

func (x *Sets) Config() Config {
	x.mu.Lock()
	defer x.mu.Unlock()
	r := x.cfg
	r.Vars = append([]Var{}, x.cfg.Vars...)
	r.Places = append([]Place{}, x.cfg.Places...)
	return r
}

func (x *Sets) SetConfig(cfg Config) {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.cfg = cfg
	x.cfg.Vars = append([]Var{}, cfg.Vars...)
	x.cfg.Places = append([]Place{}, cfg.Places...)
	if err := x.save(); err != nil {
		fmt.Println("Sets.SetConfig:", err)
	}
}

func (x *Sets) Network() Network {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.cfg.Network
}

func (x *Sets) UserConfig() settings.Config {
	x.mu.Lock()
	defer x.mu.Unlock()
	return settings.Config{
		Sections: []settings.Section{
			settings.Comport("comport", "СОМ порт", x.cfg.Comport),
			{
				Name: "chart",
				Hint: "Графики",
				Properties: []settings.Property{
					{
						Hint:         "Интервал сохранения графиков, минут",
						Name:         "save_min",
						DefaultValue: "0",
						ValueType:    settings.VtInt,
						Min:          &settings.ValueEx{Value: 0},
						Value:        strconv.Itoa(x.cfg.SaveMin),
					},
				},
			},
		},
	}
}

func (x *Sets) save() error {
	b, err := json.MarshalIndent(&x.cfg, "", "    ")
	if err != nil {
		panic(err)
	}
	return ioutil.WriteFile(configFileName(), b, 0666)
}

func (x *Sets) Save() error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.save()
}
