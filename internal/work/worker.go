package work

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/anbus/internal/data"
	"github.com/fpawel/anbus/internal/notify"
	"github.com/fpawel/anbus/internal/settings"
	"github.com/fpawel/goutils/serial/comport"
	"net"
	"strconv"
)

type worker struct {
	window          *notify.Window
	comport         *comport.Port
	sets            *anbus.Sets
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

	cfg := x.sets.Config()

	x.window.SendMsgJSON(notify.MsgUserConfig, settings.Config{
		Sections: []settings.Section{
			settings.Comport("comport", "СОМ порт", cfg.Comport),
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
						Value:        strconv.Itoa(cfg.SaveMin),
					},
				},
			},
		},
	})

	x.window.SendMsgJSON(notify.MsgNetwork, struct {
		Places []anbus.Place
		Vars   []anbus.Var
	}{
		cfg.Places,
		cfg.Vars,
	})
}
