package data

import (
	"github.com/fpawel/goutils/dbutils"
	"github.com/jmoiron/sqlx"
	"sync"
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
