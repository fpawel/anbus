package work

import (
	"github.com/fpawel/goutils/serial/comport"
	"github.com/fpawel/anbus/internal/data"
	"github.com/fpawel/anbus/internal/notify"
	"github.com/fpawel/anbus/internal/panalib"
	"github.com/fpawel/anbus/internal/settings"
	"net"
	"strconv"
	"sync"
)

type worker struct {
	window          *notify.Window
	comport         *comport.Port
	config          panalib.Config
	muConfig        sync.Mutex
	flagClose       bool
	series          *data.Series
	chModbusRequest chan modbusRequest
	ln              net.Listener
}

func (x *worker) initPeer() {
	x.window.FindPeerWindow()
	if !x.window.CheckPeerWindow() {
		return
	}

	config := x.safeGetConfig()

	x.window.SendMsgJSON(notify.MsgUserConfig, settings.Config{
		Sections: []settings.Section{
			settings.Comport("comport", "СОМ порт", config.Comport),
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
						Value:        strconv.Itoa(x.config.SaveMin),
					},
				},
			},
		},
	})

	x.window.SendMsgJSON(notify.MsgNetwork, struct {
		Places []panalib.Place
		Vars   []panalib.Var
	}{
		config.Places,
		config.Vars,
	})
}

func (x *worker) safeGetConfig() panalib.Config {
	x.muConfig.Lock()
	defer x.muConfig.Unlock()
	return x.config
}
