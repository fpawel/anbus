package work

import (
	"encoding/json"
	"github.com/fpawel/anbus/internal/settings"
	"time"
)

type CommandDataPeer uintptr

const (
	cdpSetsProperty CommandDataPeer = iota
	cdpPerformTextCommand
)

func (x *worker) onCopyData(cmd CommandDataPeer, data []byte) {
	switch cmd {

	case cdpPerformTextCommand:

		c, err := parseTxtCmd(string(data))
		if err != nil {
			x.window.SendConsoleError("%s: %v", string(data), err.Error())
			return
		}
		if c.name == "EXIT" {
			x.window.Close()
			return
		}

		if r, err := c.parseModbusRequest(); err == nil {
			x.chModbusRequest <- r
		} else {
			x.window.SendConsoleError("%s: %v", string(data), err.Error())
			return
		}

	case cdpSetsProperty:

		x.muConfig.Lock()
		defer x.muConfig.Unlock()

		var p settings.PropertyValue
		mustUnmarshalJSON(data, &p)

		switch p.Section {
		case "comport":
			switch p.Name {
			case "name":
				x.config.Comport.Serial.Name = p.Value
			case "baud":
				x.config.Comport.Serial.Baud = p.MustInt()
			case "timeout":
				x.config.Comport.Uart.ReadTimeout = time.Millisecond * p.MustDuration()
			case "timeout_byte":
				x.config.Comport.Uart.ReadTimeout = time.Millisecond * p.MustDuration()
			case "max_attempts_read":
				x.config.Comport.Uart.MaxAttemptsRead = p.MustInt()
			}
		case "chart":
			switch p.Name {
			case "save_min":
				x.config.SaveMin = p.MustInt()

			}
		}
	}
}

func mustUnmarshalJSON(data []byte, v interface{}) {
	if err := json.Unmarshal(data, v); err != nil {
		panic(err)
	}
}
