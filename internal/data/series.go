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
	db                  *sqlx.DB
	mu                  sync.Mutex
	records             []Record
	getSaveIntervalFunc func() time.Duration
}

func NewSeries(getSaveIntervalFunc func() time.Duration) *Series {
	return &Series{
		db:                  mustOpenDB(),
		getSaveIntervalFunc: getSaveIntervalFunc,
	}
}

func (x *Series) Close() error {
	return x.db.Close()
}

func (x *Series) AddRecord(addr modbus.Addr, v modbus.Var, value float64) {
	x.mu.Lock()
	defer x.mu.Unlock()

	if len(x.records) > 0 && time.Since(x.records[0].CreatedAt) > x.getSaveIntervalFunc() {
		x.save()
	}

	x.records = append(x.records, Record{
		CreatedAt:    time.Now(),
		CreatedAtStr: time.Now().Format("2006-01-02 15:04:05"),
		Addr:         addr,
		Var:          v,
		Value:        value,
	})
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

	createdAt := x.records[0].CreatedAt
	buck := x.lastBucket()

	fmt.Println("buck:", buck)
	fmt.Println("createdAt:", createdAt)
	fmt.Println("save interval:", x.getSaveIntervalFunc())
	fmt.Println("duration:", createdAt.Sub(buck.UpdatedAt))

	if buck.BucketID == 0 || createdAt.Sub(buck.UpdatedAt) > x.getSaveIntervalFunc() {
		x.db.MustExec(`INSERT INTO bucket DEFAULT VALUES;`)
		buck = x.lastBucket()
	}
	queryStr := `INSERT INTO series(bucket_id, addr, var, value, stored_at)  VALUES `
	for i, a := range x.records {

		s := fmt.Sprintf("(%d, %d, %d, %v, julianday('%s'))", buck.BucketID,
			a.Addr, a.Var, a.Value,
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

type bucket struct {
	BucketID  int64      `db:"bucket_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	Year      int        `db:"year"`
	Month     time.Month `db:"month"`
	Day       int        `db:"day"`
}

func (x *Series) lastBucket() bucket {
	var xs []bucket
	dbutils.MustSelect(x.db, &xs,
		`
SELECT bucket_id, 
       created_at, 
       updated_at         
FROM last_bucket;`)
	if len(xs) == 0 {
		return bucket{}
	}
	xs[0].CreatedAt.Add(time.Hour * 3)
	return xs[0]
}
