package main

import (
	"flag"
	"github.com/fpawel/anbus/internal/app"
	"github.com/fpawel/gohelp/must"
	"github.com/lxn/win"
	"github.com/powerman/structlog"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defaultLogLevelStr := os.Getenv("ANBUS_LOG_LEVEL")
	if len(strings.TrimSpace(defaultLogLevelStr)) == 0 {
		defaultLogLevelStr = "info"
	}

	hideCon := flag.Bool("hide.con", false, "hide console window")
	logLevel := flag.String("log.level", defaultLogLevelStr, "log `level` (debug|info|warn|err)")

	flag.Parse()

	if *hideCon {
		win.ShowWindow(win.GetConsoleWindow(), win.SW_HIDE)
	}

	structlog.DefaultLogger.
		SetPrefixKeys(
			structlog.KeyApp, structlog.KeyPID, structlog.KeyLevel, structlog.KeyUnit, structlog.KeyTime,
		).
		SetDefaultKeyvals(
			structlog.KeyApp, filepath.Base(os.Args[0]),
			structlog.KeySource, structlog.Auto,
		).
		SetSuffixKeys(
			structlog.KeyStack,
		).
		SetSuffixKeys(structlog.KeySource).
		SetKeysFormat(map[string]string{
			structlog.KeyTime:   " %[2]s",
			structlog.KeySource: " %6[2]s",
			structlog.KeyUnit:   " %6[2]s",
		}).SetTimeFormat("15:04:05").
		SetLogLevel(structlog.ParseLevel(*logLevel))

	must.AbortIf = must.PanicIf
	app.Run()
}
