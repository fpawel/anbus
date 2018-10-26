package settings

import (
	"github.com/fpawel/goutils/serial/comport"
	"strconv"
	"time"
)

func Comport(name, hint string, c comport.Config) Section {

	return Section{
		Name: name,
		Hint: hint,
		Properties: []Property{
			{
				Name:         "name",
				Hint:         "Имя СОМ порта",
				ValueType:    VtComportName,
				DefaultValue: "COM1",
				Value:        c.Serial.Name,
			},
			{
				Name:         "baud",
				Hint:         "Скорость приёмопередачи, бод",
				ValueType:    VtBaud,
				DefaultValue: "9600",
				Value:        strconv.Itoa(c.Serial.Baud),
			},
			{
				Name:         "timeout",
				Hint:         "Таймаут посылки, мс",
				ValueType:    VtInt,
				Min:          &ValueEx{10},
				Max:          &ValueEx{10000},
				DefaultValue: "1000",
				Value:        timeMillis(c.Uart.ReadTimeout),
			},
			{
				Name:         "timeout_byte",
				Hint:         "Таймаут байта, мс",
				ValueType:    VtInt,
				Min:          &ValueEx{10},
				Max:          &ValueEx{100},
				DefaultValue: "50",
				Value:        timeMillis(c.Uart.ReadByteTimeout),
			},
			{
				Name:         "max_attempts_read",
				Hint:         "Макс. кол-во попыток",
				ValueType:    VtInt,
				Min:          &ValueEx{0},
				Max:          &ValueEx{10},
				DefaultValue: "0",
				Value:        strconv.Itoa(c.Uart.MaxAttemptsRead),
			},
		},
	}

}

func timeMillis(t time.Duration) string {
	return strconv.Itoa(int(t / time.Millisecond))
}
