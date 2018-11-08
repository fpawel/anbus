package data

import (
	"fmt"
	"github.com/fpawel/goutils/serial/modbus"
)

type DeletePointsRequest struct {
	Addr     modbus.Addr
	Var      modbus.Var
	BucketID int64
	ValueMinimum,
	ValueMaximum float64
	TimeMinimum,
	TimeMaximum Time
}

func (x DeletePointsRequest) String() string {
	return fmt.Sprintf(`
DELETE FROM series 
WHERE bucket_id = %d AND 
      addr = %d AND 
      var = %d AND  
      value >= %v AND 
      value <= %v AND 
      stored_at >= julianday('%v') AND 
      stored_at <= julianday('%v');`,

		x.BucketID, x.Addr, x.Var,
		x.ValueMinimum, x.ValueMaximum,
		x.TimeMinimum.Time().Format(timeFormat),
		x.TimeMaximum.Time().Format(timeFormat),
	)
}
