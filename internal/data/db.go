package data

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/goutils/dbutils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

//go:generate go run github.com/fpawel/goutils/dbutils/sqlstr/...

func mustOpenDB() *sqlx.DB {
	db := dbutils.MustOpen(anbus.AppName.DataFileName("series.sqlite"), "sqlite3")
	db.MustExec(SQLCreate)
	return db
}

func lastBucket(db *sqlx.DB) Bucket {
	var xs []Bucket
	dbutils.MustSelect(db, &xs,
		`SELECT bucket_id, created_at, updated_at FROM last_bucket;`)
	if len(xs) == 0 {
		return Bucket{}
	}
	xs[0].CreatedAt.Add(time.Hour * 3)
	return xs[0]
}
