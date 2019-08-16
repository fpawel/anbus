package cfg

import (
	"encoding/json"
	"fmt"
	"github.com/fpawel/comm"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gohelp/winapp"
	"github.com/powerman/must"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	ConfigNetwork
	ConfigEditable
}

type ConfigEditable struct {
	Comm        comm.Config `toml:"comm" comment:"параметры приёмопередачи"`
	ComportName string      `toml:"comport_name" comment:"имя СОМ порта"`
	ComportBaud int         `toml:"comport_baud" comment:"скорость передачи, бод"`
	SaveSeries  bool        `toml:"save_series" comment:"сохранять графики"`
}

type ConfigNetwork struct {
	Places []Place
	Vars   []DevVar
}

type Place struct {
	Addr  modbus.Addr
	Check bool
}

type DevVar struct {
	Code  modbus.Var
	Check bool
}

type Node struct {
	Place    int
	Addr     modbus.Addr
	VarCode  modbus.Var
	VarIndex int
}

func (x ConfigNetwork) ToggleChecked() {
	f := len(x.Nodes()) > 0
	for i := range x.Places {
		x.Places[i].Check = !f
	}
	for i := range x.Vars {
		x.Vars[i].Check = !f
	}
}

func (x ConfigNetwork) NextNode(n Node) Node {
	xs := x.Nodes()
	if len(xs) == 0 {
		return Node{Place: -1}
	}
	for i, vb := range xs {
		if vb == n && i < len(xs)-1 {
			return xs[i+1]
		}
	}
	return xs[0]
}

func (x ConfigNetwork) Nodes() (xs []Node) {
	for place, a := range x.Places {
		if !a.Check {
			continue
		}
		for varIndex, v := range x.Vars {
			if !v.Check {
				continue
			}
			xs = append(xs, Node{place, a.Addr, v.Code, varIndex,})
		}
	}
	return
}

func (x Config) Save() error {
	b, err := json.MarshalIndent(x, "", "    ")
	if err != nil {
		return err
	}
	configFileName, err := winapp.ProfileFileName(".anbus", "cfg.json")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFileName, b, 0666)
}

func Get() (result Config) {
	mu.Lock()
	defer mu.Unlock()
	must.UnmarshalJSON(must.MarshalJSON(cfg), &result)
	return result
}

func Set(c Config) {
	mu.Lock()
	defer mu.Unlock()
	must.UnmarshalJSON(must.MarshalJSON(c), &cfg)
	save()
}

func save() {
	b, err := json.MarshalIndent(&cfg, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(fileName(), b, 0666); err != nil {
		panic(err)
	}
}

func fileName() string {
	return filepath.Join(filepath.Dir(os.Args[0]), "anbus.cfg.json")
}

var (
	cfg = func() (c Config) {
		c = Config{
			ConfigEditable: ConfigEditable{
				ComportName: "COM1",
				ComportBaud: 9600,
				Comm: comm.Config{
					MaxAttemptsRead:       2,
					ReadByteTimeoutMillis: 30,
					ReadTimeoutMillis:     500,
				},
				SaveSeries: true,
			},
			ConfigNetwork: ConfigNetwork{
				Places: []Place{{Addr: 1, Check: true}},
				Vars: []DevVar{{Code: 0, Check: false,
				}},
			},
		}
		b, err := ioutil.ReadFile(fileName())
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = json.Unmarshal(b, &c); err == nil {
			fmt.Println(err)
			return
		}
		return
	}()

	mu sync.Mutex
)
