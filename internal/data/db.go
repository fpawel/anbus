package data

import (
	"github.com/fpawel/anbus/internal/anbus"
	"github.com/fpawel/goutils/dbutils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:generate go run github.com/fpawel/goutils/dbutils/sqlstr/...

func mustOpenDB() *sqlx.DB {
	db := dbutils.MustOpen(anbus.AppName.DataFileName("series.sqlite"), "sqlite3")
	db.MustExec(SQLCreate)
	return db
}
