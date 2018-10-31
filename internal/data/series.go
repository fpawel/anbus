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
	records []record
}

type record struct {
	storedAt time.Time
	nVar     modbus.Var
	addr     modbus.Addr
	value    float64
}

func NewSeries() *Series {
	return &Series{
		db: mustOpenDB(),
	}
}

func (x *Series) Close() error {
	return x.db.Close()
}

func (x *Series) AddRecord(addr modbus.Addr, v modbus.Var, value float64) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.records = append(x.records, record{
		storedAt: time.Now(),
		addr:     addr,
		nVar:     v,
		value:    value,
	})
	if time.Since(x.records[0].storedAt) > time.Minute {
		x.save()
	}
}

func (x *Series) Save() {

	x.mu.Lock()
	defer x.mu.Unlock()
	x.save()
}

func (x *Series) Buckets() *BucketsSvc {
	return &BucketsSvc{
		db: x.db,
		mu: &x.mu,
	}
}

func (x *Series) save() {
	if len(x.records) == 0 {
		return
	}

	buck := x.lastBucket()

	if buck.BucketID == 0 || x.records[0].storedAt.Sub(buck.UpdatedAt) > time.Minute {
		x.db.MustExec(`INSERT INTO bucket DEFAULT VALUES;`)
		buck = x.lastBucket()
	}
	queryStr := `INSERT INTO series(bucket_id, Addr, var, Value, stored_at)  VALUES `
	for i, a := range x.records {

		s := fmt.Sprintf("(%d, %d, %d, %v, julianday('%s'))", buck.BucketID,
			a.addr, a.nVar, a.value,
			a.storedAt.Format("2006-01-02 15:04:05.000"))
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

func (x *Series) lastBucket() Bucket {
	var xs []Bucket
	dbutils.MustSelect(x.db, &xs,
		`SELECT bucket_id, created_at, updated_at FROM last_bucket;`)
	if len(xs) == 0 {
		return Bucket{}
	}
	xs[0].CreatedAt.Add(time.Hour * 3)
	return xs[0]
}
