package ser

import (
	"fmt"
	"github.com/fpawel/anbus/internal/data"
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
	StoredAt time.Time
	Var      modbus.Var
	Addr     modbus.Addr
	Value    float64
}

func NewSeries() *Series {
	return &Series{
		db: data.MustOpenDB(),
	}
}

func OpenFile(fileName string) (*Series, error) {

	db := sqlx.MustConnect(data.DriverName, fileName)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := db.Close(); err != nil {
		return nil, err
	}
	db = sqlx.MustConnect(data.DriverName, fileName)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Series{db: db}, nil
}

func (x *Series) Close() error {
	return x.db.Close()
}

func (x *Series) AddRecord(addr modbus.Addr, v modbus.Var, value float64) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.records = append(x.records, record{
		StoredAt: time.Now(),
		Addr:     addr,
		Var:      v,
		Value:    value,
	})
	if time.Since(x.records[0].StoredAt) > time.Minute {
		x.save()
	}
}

// save - залочить базу данных, сохранить точки из кеша, очистить кеш, разлочить базу данных
// Вернуть true, если в таблицу bucket была добавлена новая запись
func (x *Series) Save() bool {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.save()
}

// save - сохранить точки из кеша, очистить кеш.
// Вернуть true, если в таблицу bucket была добавлена новая запись
func (x *Series) save() bool {
	if len(x.records) == 0 {
		return false
	}
	buck := x.lastBucket()
	result := buck.BucketID == 0 || x.records[0].StoredAt.Sub(buck.UpdatedAt) > time.Minute
	if result {
		x.db.MustExec(`INSERT INTO bucket DEFAULT VALUES;`)
		buck = x.lastBucket()
	}
	x.doSave(buck.BucketID)
	return result
}

func (x *Series) doSave(bucketID int64) {
	queryStr := `INSERT INTO series(bucket_id, Addr, var, Value, stored_at)  VALUES `
	for i, a := range x.records {

		s := fmt.Sprintf("(%d, %d, %d, %v, julianday('%s'))", bucketID,
			a.Addr, a.Var, a.Value,
			a.StoredAt.Format("2006-01-02 15:04:05.000"))
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

func (x *Series) lastBucket() data.Bucket {
	return data.LastBucket(x.db)
}

func (x *Series) Years() []int {
	x.mu.Lock()
	defer x.mu.Unlock()
	return data.GetYears(x.db)
}

func (x *Series) MonthsOfYear(y int) []int {
	x.mu.Lock()
	defer x.mu.Unlock()
	return data.GetMonthsOfYear(x.db, y)
}

func (x *Series) DaysOfYearMonth(y, m int) []int {
	x.mu.Lock()
	defer x.mu.Unlock()
	return data.GetDaysOfYearMonth(x.db, y, m)
}

func (x *Series) BucketsOfDayYearMonth(y, m, d int) []data.Bucket {
	x.mu.Lock()
	defer x.mu.Unlock()
	return data.GetBuckets(x.db, y, m, d)
}

func (x *Series) VarsByBucketID(bucketID int64) []int {
	x.mu.Lock()
	defer x.mu.Unlock()
	return data.GetVarsByBucketID(x.db, bucketID)
}

func (x *Series) AddressesByBucketID(bucketID int64) []int {
	x.mu.Lock()
	defer x.mu.Unlock()
	return data.GetAddressesByBucketID(x.db, bucketID)
}

func (x *Series) DeletePoints(request data.DeletePointsRequest) int64 {

	x.mu.Lock()
	defer x.mu.Unlock()

	if request.BucketID == 0 {
		request.BucketID = x.lastBucket().BucketID
		var records []record
		for _, a := range x.records {
			f := a.Addr == request.Addr && a.Var == request.Var &&
				a.StoredAt.After(request.TimeMinimum.Time()) && a.StoredAt.Before(request.TimeMaximum.Time()) &&
				a.Value >= request.ValueMinimum &&
				a.Value <= request.ValueMaximum
			if !f {
				records = append(records, a)
			}
		}
	}
	return data.DeletePoints(x.db, request)
}

func (x *Series) PointsByBucketID(bucketID int64) []data.ChartPoint {
	x.mu.Lock()
	defer x.mu.Unlock()
	return data.GetPoints(x.db, bucketID)
}
