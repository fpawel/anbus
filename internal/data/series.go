package data

import (
	"fmt"
	"github.com/fpawel/goutils/dbutils"
	"github.com/fpawel/goutils/serial/modbus"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type Series struct {
	db      *sqlx.DB
	mu      sync.Mutex
	records []Record
}

func NewSeries() *Series {
	return &Series{
		db: mustOpenDB(),
	}
}

func (x *Series) Close() error {
	return x.db.Close()
}

//func (x *Series) last() Record {
//	if len(x.records) > 0 {
//		return x.records[len(x.records)-1]
//	}
//	return Record{}
//}

func (x *Series) AddRecord(addr modbus.Addr, v modbus.Var, value float64, saveInterval time.Duration) {
	if len(x.records) > 0 && time.Since(x.records[0].CreatedAt) > saveInterval {
		x.Upload(saveInterval)
	}
	x.records = append(x.records, Record{
		CreatedAt:    time.Now(),
		CreatedAtStr: time.Now().Format("2006-01-02 15:04:05"),
		Addr:         addr,
		Var:          v,
		Value:        value,
	})
}

func (x *Series) lastBucket() (int64, time.Time) {
	var xs []struct {
		BucketID int64      `db:"bucket_id"`
		Year     int        `db:"year"`
		Month    time.Month `db:"month"`
		Day      int        `db:"day"`
		Hour     int        `db:"hour"`
		Minute   int        `db:"minute"`
		Second   int        `db:"second"`
		StoredAt string     `db:"stored_at"`
	}
	dbutils.MustSelect(x.db, &xs, `SELECT * FROM last_value;`)
	if len(xs) == 0 {
		return 0, time.Time{}
	}
	a := xs[0]
	return a.BucketID, time.Date(a.Year, a.Month, a.Day, a.Hour, a.Minute, a.Second, 0, time.UTC)
}

func (x *Series) Upload(saveInterval time.Duration) {
	if len(x.records) == 0 {
		return
	}

	x.mu.Lock()
	defer x.mu.Unlock()

	createdAt := x.records[0].CreatedAt
	bucketID, storedAt := x.lastBucket()

	if bucketID == 0 || createdAt.Sub(storedAt) > saveInterval {
		r := x.db.MustExec(`INSERT INTO bucket (created_at) VALUES (?);`, createdAt)
		var err error
		bucketID, err = r.LastInsertId()
		if err != nil {
			panic(err)
		}
	}
	queryStr := `INSERT INTO series(bucket_id, addr, var, value, stored_at)  VALUES `
	for i, a := range x.records {

		s := fmt.Sprintf("(%d, %d, %d, %v, julianday('%s'))", bucketID, a.Addr, a.Var, a.Value,
			a.CreatedAt.Format("2006-01-02 15:04:05.000"))
		if i == len(x.records)-1 {
			s += ";"
		} else {
			s += ", "
		}
		queryStr += s
	}
	x.db.MustExec(queryStr)
	x.records = nil
}
