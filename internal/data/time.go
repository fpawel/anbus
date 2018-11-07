package data

import "time"

type Time struct {
	Year                                   int
	Month                                  time.Month
	Day, Hour, Minute, Second, Millisecond int
}

func (x Time) Time() time.Time {
	return time.Date(
		x.Year, x.Month, x.Day,
		x.Hour, x.Minute, x.Second,
		x.Millisecond*int(time.Millisecond/time.Nanosecond),
		time.Local)
}
