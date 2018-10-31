package data

import (
	"github.com/fpawel/goutils/dbutils"
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type BucketsSvc struct {
	db *sqlx.DB
	mu *sync.Mutex
}

func (x *BucketsSvc) Years(_ struct{}, years *[]int) error {
	dbutils.MustSelect(x.db, years, `SELECT DISTINCT year FROM bucket_time;`)
	return nil
}

func (x *BucketsSvc) Months(y [1]int, months *[]int) error {
	dbutils.MustSelect(x.db, months,
		`SELECT DISTINCT month FROM bucket_time WHERE year = ?;`, y[0])
	return nil
}

func (x *BucketsSvc) Days(p [2]int, days *[]int) error {
	dbutils.MustSelect(x.db, days,
		`
SELECT DISTINCT day 
FROM bucket_time 
WHERE year = ? AND month = ?;`,
		p[0], p[1])
	return nil
}

func (x *BucketsSvc) Buckets(p [3]int, buckets *[]Bucket) error {
	dbutils.MustSelect(x.db, buckets,
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
	dbutils.MustSelect(x.db, vars,
		`SELECT DISTINCT var FROM series WHERE bucket_id = ?;`, p[0])
	return nil
}

func (x *BucketsSvc) Records(p [1]int, r *[][10]float64) error {

	var xs []struct {
		StoredAt string      `db:"stored_at"`
		Var      modbus.Var  `db:"var"`
		Addr     modbus.Addr `db:"addr"`
		Value    float64     `db:"value"`
	}

	dbutils.MustSelect(x.db, &xs,
		`
SELECT addr, var, value, strftime('%Y-%m-%d %H:%M:%f', stored_at) AS stored_at
FROM series 
WHERE bucket_id = ?;`, p[0])

	for _, v := range xs {
		t, err := time.Parse("2006-01-02 15:04:05.000", v.StoredAt)
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
