package settings

import (
	"strconv"
	"time"
)

type PropertyValue struct {
	Value, Section, Name string
}

func (x PropertyValue) MustInt() int {
	n, err := strconv.Atoi(x.Value)
	if err != nil {
		panic(err)
	}
	return n
}

func (x PropertyValue) MustDuration() time.Duration {
	return time.Duration(x.MustInt())
}
