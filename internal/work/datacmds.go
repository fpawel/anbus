package work

import (
	"fmt"
)

type CommandDataPeer uintptr

const (
	cdpPerformTextCommand CommandDataPeer = iota
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
			if err := x.window.Close(); err != nil {
				fmt.Println("close window: unexpected error:", err)
			}
			return
		}

		if r, err := c.parseRequest(); err == nil {
			x.chModbusRequest <- r
		} else {
			x.window.SendConsoleError("%s: %v", string(data), err.Error())
			return
		}
	}
}
