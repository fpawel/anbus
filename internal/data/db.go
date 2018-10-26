package data

import (
	"github.com/fpawel/goutils/dbutils"
	"github.com/fpawel/anbus/internal/panalib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:generate go run github.com/fpawel/goutils/dbutils/sqlstr/...

func mustOpenDB() *sqlx.DB {
	return dbutils.MustOpen(panalib.AppName.DataFileName("series.sqlite"), "sqlite3")
}
