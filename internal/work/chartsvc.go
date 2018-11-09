package work

import (
	"fmt"
	"github.com/fpawel/anbus/internal/chart"
	"github.com/fpawel/anbus/internal/data"
)

type ChartSvc struct {
	series *chart.Series
}

func (x *ChartSvc) Years(_ struct{}, years *[]int) error {
	*years = x.series.Years()
	return nil
}

func (x *ChartSvc) Months(r struct{ Year int }, months *[]int) error {
	*months = x.series.MonthsOfYear(r.Year)
	return nil
}

func (x *ChartSvc) Days(r struct{ Year, Month int }, days *[]int) error {
	*days = x.series.DaysOfYearMonth(r.Year, r.Month)
	return nil
}

func (x *ChartSvc) Buckets(r struct{ Year, Month, Day int }, buckets *[]data.Bucket) error {
	*buckets = x.series.BucketsOfDayYearMonth(r.Year, r.Month, r.Day)
	return nil
}

func (x *ChartSvc) Vars(r struct{ BucketID int64 }, vars *[]int) error {
	*vars = x.series.VarsByBucketID(r.BucketID)
	return nil
}

func (x *ChartSvc) Addresses(r struct{ BucketID int64 }, addresses *[]int) error {
	*addresses = x.series.AddressesByBucketID(r.BucketID)
	return nil
}

func (x *ChartSvc) DeletePoints(r struct{ data.DeletePointsRequest }, result *int64) error {

	fmt.Println(r.DeletePointsRequest)
	*result = x.series.DeletePoints(r.DeletePointsRequest)
	return nil
}

func (x *ChartSvc) Points(r struct{ BucketID int64 }, points *[]data.ChartPoint) error {
	*points = x.series.PointsByBucketID(r.BucketID)
	return nil
}
