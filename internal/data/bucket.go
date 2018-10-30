package data

import "time"

type Bucket struct {
	BucketID  int64      `db:"bucket_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	Year      int        `db:"year"`
	Month     time.Month `db:"month"`
	Day       int        `db:"day"`
}
