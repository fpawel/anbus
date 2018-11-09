package data

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/goutils/dbutils"
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

//go:generate go run github.com/fpawel/goutils/dbutils/sqlstr/...

type ChartPoint = [10]float64

const DriverName = "sqlite3"

func MustOpenDB() *sqlx.DB {
	db := dbutils.MustOpen(anbus.DataFileName(), DriverName)
	db.MustExec(SQLCreate)
	return db
}

func GetYears(db *sqlx.DB) (years []int) {
	dbutils.MustSelect(db, &years, `SELECT DISTINCT year FROM bucket_time;`)
	return
}

func GetMonthsOfYear(db *sqlx.DB, y int) (months []int) {
	dbutils.MustSelect(db, &months,
		`SELECT DISTINCT month FROM bucket_time WHERE year = ?;`, y)
	return
}

func GetDaysOfYearMonth(db *sqlx.DB, y, m int) (days []int) {
	dbutils.MustSelect(db, &days,
		`
SELECT DISTINCT day 
FROM bucket_time 
WHERE year = ? AND month = ?;`,
		y, m)
	return
}

func GetBuckets(db *sqlx.DB, year, month, day int) (buckets []Bucket) {
	dbutils.MustSelect(db, &buckets,
		`
SELECT * FROM bucket_time 
WHERE year = ? AND month = ? AND day = ?;`,
		year, month, day)

	for i := range buckets {
		(buckets)[i].CreatedAt = (buckets)[i].CreatedAt.Add(3 * time.Hour)
		(buckets)[i].UpdatedAt = (buckets)[i].UpdatedAt.Add(3 * time.Hour)
	}
	return
}

func LastBucket(db *sqlx.DB) Bucket {
	var xs []Bucket
	dbutils.MustSelect(db, &xs,
		`SELECT bucket_id, created_at, updated_at FROM last_bucket;`)
	if len(xs) == 0 {
		return Bucket{}
	}
	xs[0].CreatedAt.Add(time.Hour * 3)
	return xs[0]
}

func GetVarsByBucketID(db *sqlx.DB, bucketID int64) (vars []int) {
	dbutils.MustSelect(db, &vars,
		`SELECT DISTINCT var FROM series WHERE bucket_id = ?;`, bucketID)
	return
}

func GetAddressesByBucketID(db *sqlx.DB, bucketID int64) (addresses []int) {
	dbutils.MustSelect(db, &addresses,
		`SELECT DISTINCT addr FROM series WHERE bucket_id = ?;`, bucketID)
	return
}

const timeFormat = "2006-01-02 15:04:05.000"

func DeletePoints(db *sqlx.DB, request DeletePointsRequest) int64 {
	r := db.MustExec(
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
	return n
}

func GetPoints(db *sqlx.DB, bucketID int64) (points []ChartPoint) {

	var xs []struct {
		StoredAt string      `db:"stored_at"`
		Var      modbus.Var  `db:"var"`
		Addr     modbus.Addr `db:"addr"`
		Value    float64     `db:"value"`
	}

	dbutils.MustSelect(db, &xs,
		`
SELECT addr, var, value, strftime('%Y-%m-%d %H:%M:%f', stored_at) AS stored_at
FROM series 
WHERE bucket_id = ?;`, bucketID)

	for _, v := range xs {
		t, err := time.Parse(timeFormat, v.StoredAt)
		if err != nil {
			panic(err)
		}
		points = append(points, ChartPoint{
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
	return
}
