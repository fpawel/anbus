package svc

import (
	"fmt"
	"github.com/fpawel/anbus/internal/data"
	"github.com/fpawel/anbus/internal/data/ser"
)

type ChartSvc struct {
	m map[string]*ser.Series
}

type ChartSvcRequest struct {
	FileName string
}

func NewChartSvc() (*ChartSvc, *ser.Series) {
	series := ser.NewSeries()
	return &ChartSvc{
		m: map[string]*ser.Series{"": series},
	}, series
}

func (x *ChartSvc) series(fileName string) (*ser.Series, error) {
	if v, ok := x.m[fileName]; ok {
		return v, nil
	}
	series, err := ser.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	x.m[fileName] = series
	return series, nil
}

func (x *ChartSvc) Years(r ChartSvcRequest, years *[]int) error {
	p, err := x.series(r.FileName)
	if err != nil {
		return err
	}
	*years = p.Years()
	return nil
}

func (x *ChartSvc) Months(r struct {
	ChartSvcRequest
	Year int
}, months *[]int) error {
	p, err := x.series(r.FileName)
	if err != nil {
		return err
	}
	*months = p.MonthsOfYear(r.Year)
	return nil
}

func (x *ChartSvc) Days(r struct {
	ChartSvcRequest
	Year, Month int
}, days *[]int) error {
	p, err := x.series(r.FileName)
	if err != nil {
		return err
	}
	*days = p.DaysOfYearMonth(r.Year, r.Month)
	return nil
}

func (x *ChartSvc) Buckets(r struct {
	ChartSvcRequest
	Year, Month, Day int
}, buckets *[]data.Bucket) error {
	p, err := x.series(r.FileName)
	if err != nil {
		return err
	}
	*buckets = p.BucketsOfDayYearMonth(r.Year, r.Month, r.Day)
	return nil
}

func (x *ChartSvc) Vars(r struct {
	ChartSvcRequest
	BucketID int64
}, vars *[]int) error {
	p, err := x.series(r.FileName)
	if err != nil {
		return err
	}
	*vars = p.VarsByBucketID(r.BucketID)
	return nil
}

func (x *ChartSvc) Addresses(r struct {
	ChartSvcRequest
	BucketID int64
}, addresses *[]int) error {
	p, err := x.series(r.FileName)
	if err != nil {
		return err
	}
	*addresses = p.AddressesByBucketID(r.BucketID)
	return nil
}

func (x *ChartSvc) DeletePoints(r struct {
	ChartSvcRequest
	data.DeletePointsRequest
}, result *int64) error {

	p, err := x.series(r.FileName)
	if err != nil {
		return err
	}
	fmt.Println(r.DeletePointsRequest)
	*result = p.DeletePoints(r.DeletePointsRequest)
	return nil
}

func (x *ChartSvc) Points(r struct {
	ChartSvcRequest
	BucketID int64
}, points *[]data.ChartPoint) error {
	p, err := x.series(r.FileName)
	if err != nil {
		return err
	}
	*points = p.PointsByBucketID(r.BucketID)
	return nil
}
