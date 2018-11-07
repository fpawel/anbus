package data

import (
	"fmt"
	"github.com/fpawel/goutils/dbutils"
	"github.com/fpawel/goutils/serial/modbus"
	"time"
)

type BucketsSvc struct {
	p *Series
}

const timeFormat = "2006-01-02 15:04:05.000"

func (x *BucketsSvc) Years(_ struct{}, years *[]int) error {
	x.p.mu.Lock()
	defer x.p.mu.Unlock()
	dbutils.MustSelect(x.p.db, years, `SELECT DISTINCT year FROM bucket_time;`)
	return nil
}

func (x *BucketsSvc) Months(y [1]int, months *[]int) error {
	x.p.mu.Lock()
	defer x.p.mu.Unlock()
	dbutils.MustSelect(x.p.db, months,
		`SELECT DISTINCT month FROM bucket_time WHERE year = ?;`, y[0])
	return nil
}

func (x *BucketsSvc) Days(p [2]int, days *[]int) error {
	x.p.mu.Lock()
	defer x.p.mu.Unlock()
	dbutils.MustSelect(x.p.db, days,
		`
SELECT DISTINCT day 
FROM bucket_time 
WHERE year = ? AND month = ?;`,
		p[0], p[1])
	return nil
}

func (x *BucketsSvc) Buckets(p [3]int, buckets *[]Bucket) error {
	x.p.mu.Lock()
	defer x.p.mu.Unlock()
	dbutils.MustSelect(x.p.db, buckets,
		`
SELECT * FROM bucket_time 
WHERE year = ? AND month = ? AND day = ?;`,
		p[0], p[1], p[2])

	for i := range *buckets {
		(*buckets)[i].CreatedAt = (*buckets)[i].CreatedAt.Add(3 * time.Hour)
		(*buckets)[i].UpdatedAt = (*buckets)[i].UpdatedAt.Add(3 * time.Hour)
	}
	return nil
}

func (x *BucketsSvc) Vars(p [1]int, vars *[]int) error {
	x.p.mu.Lock()
	defer x.p.mu.Unlock()
	dbutils.MustSelect(x.p.db, vars,
		`SELECT DISTINCT var FROM series WHERE bucket_id = ?;`, p[0])
	return nil
}

func (x *BucketsSvc) Addresses(p [1]int, vars *[]int) error {
	x.p.mu.Lock()
	defer x.p.mu.Unlock()
	dbutils.MustSelect(x.p.db, vars,
		`SELECT DISTINCT addr FROM series WHERE bucket_id = ?;`, p[0])
	return nil
}

type DeletePointsRequest struct {
	Addr                       modbus.Addr
	Var                        modbus.Var
	BucketID                   int64
	ValueMinimum, ValueMaximum float64
	TimeMinimum, TimeMaximum   Time
}

func (x *BucketsSvc) DeletePoints(request DeletePointsRequest, result *int64) error {

	x.p.mu.Lock()
	defer x.p.mu.Unlock()

	if request.BucketID == 0 {
		request.BucketID = x.p.lastBucket().BucketID
		var records []record
		for _, a := range x.p.records {
			f := a.Addr == request.Addr && a.Var == request.Var &&
				a.StoredAt.After(request.TimeMinimum.Time()) && a.StoredAt.Before(request.TimeMaximum.Time()) &&
				a.Value >= request.ValueMinimum &&
				a.Value <= request.ValueMaximum
			if !f {
				records = append(records, a)
			}
		}
	}
	fmt.Printf(`
DELETE FROM series 
WHERE bucket_id = %d AND 
      addr = %d AND 
      var = %d AND  
      value >= %v AND 
      value <= %v AND 
      stored_at >= julianday('%v') AND 
      stored_at <= julianday('%v');`,

		request.BucketID, request.Addr, request.Var,
		request.ValueMinimum, request.ValueMaximum,
		request.TimeMinimum.Time().Format(timeFormat),
		request.TimeMaximum.Time().Format(timeFormat),
	)

	r := x.p.db.MustExec(
		`
DELETE FROM series 
WHERE bucket_id = ? AND 
      addr = ? AND 
      var = ? AND  
      value >= ? AND 
      value <= ? AND 
      stored_at >= julianday(?) AND 
      stored_at <= julianday(?);`, request.BucketID, request.Addr, request.Var,
		request.ValueMinimum, request.ValueMaximum,
		request.TimeMinimum.Time().Format(timeFormat),
		request.TimeMaximum.Time().Format(timeFormat))

	n, err := r.RowsAffected()
	if err != nil {
		panic(err)
	}
	*result = n

	return nil
}

func (x *BucketsSvc) Records(p [1]int, r *[][10]float64) error {
	x.p.mu.Lock()
	defer x.p.mu.Unlock()

	var xs []struct {
		StoredAt string      `db:"stored_at"`
		Var      modbus.Var  `db:"var"`
		Addr     modbus.Addr `db:"addr"`
		Value    float64     `db:"value"`
	}

	dbutils.MustSelect(x.p.db, &xs,
		`
SELECT addr, var, value, strftime('%Y-%m-%d %H:%M:%f', stored_at) AS stored_at
FROM series 
WHERE bucket_id = ?;`, p[0])

	for _, v := range xs {
		t, err := time.Parse(timeFormat, v.StoredAt)
		if err != nil {
			panic(err)
		}
		*r = append(*r, [10]float64{
			float64(v.Addr),
			float64(v.Var),
			float64(t.Year()),
			float64(t.Month()),
			float64(t.Day()),
			float64(t.Hour()),
			float64(t.Minute()),
			float64(t.Second()),
			float64(t.Nanosecond() / int(time.Millisecond)),
			v.Value,
		})

	}

	return nil
}
