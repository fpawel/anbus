package main

import (
	"github.com/fpawel/anbus/internal/api"
	"github.com/fpawel/anbus/internal/api/types"
	"github.com/fpawel/dseries"
	"github.com/fpawel/gohelp/delphi/delphirpc"
	"os"
	"path/filepath"
	r "reflect"
)

func main() {

	delphirpc.WriteSources(delphirpc.SrcServices{
		Dir: filepath.Join(os.Getenv("DELPHIPATH"),
			"src", "github.com", "fpawel", "anbusgui", "api"),
		Types: []r.Type{
			r.TypeOf((*api.ConfigSvc)(nil)),
			r.TypeOf((*dseries.ChartsSvc)(nil)),
			r.TypeOf((*api.TaskSvc)(nil)),
		},
	}, delphirpc.SrcNotify{
		Dir: filepath.Join(os.Getenv("GOPATH"),
			"src", "github.com", "fpawel", "anbus", "internal", "api", "notify"),
		Types: []delphirpc.NotifyServiceType{

			{
				"WriteConsoleInfo",
				r.TypeOf((*string)(nil)).Elem(),
			},
			{
				"WriteConsoleError",
				r.TypeOf((*string)(nil)).Elem(),
			},
			{
				"Status",
				r.TypeOf((*string)(nil)).Elem(),
			},
			{
				"WorkError",
				r.TypeOf((*string)(nil)).Elem(),
			},
			{
				"ReadVar",
				r.TypeOf((*types.ReadVar)(nil)).Elem(),
			},
			{
				"NewSeries",
				r.TypeOf((*struct{})(nil)).Elem(),
			},
		},
		PeerWindowClassName:   "TAnbusMainForm",
		ServerWindowClassName: "AnbusServerWindow",
	})

}
