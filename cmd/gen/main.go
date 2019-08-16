package main

import (
	"github.com/fpawel/anbus/internal/api"
	"github.com/fpawel/dseries"
	"github.com/fpawel/gohelp/delphi/delphirpc"
	"os"
	"path/filepath"
	r "reflect"
)

func main() {

	delphirpc.WriteSources(delphirpc.SrcServices{
		Dir: filepath.Join(os.Getenv("DELPHIPATH"),
			"src", "github.com", "fpawel", "anbus", "api"),
		Types: []r.Type{
			r.TypeOf((*api.ConfigSvc)(nil)),
			r.TypeOf((*dseries.ChartsSvc)(nil)),
		},
	}, delphirpc.SrcNotify{
		Dir: filepath.Join(os.Getenv("GOPATH"),
			"src", "github.com", "fpawel", "anbus", "internal", "api", "notify"),
		Types: []delphirpc.NotifyServiceType{
			{
				"Info",
				r.TypeOf((*string)(nil)).Elem(),
			},
			{
				"Error",
				r.TypeOf((*string)(nil)).Elem(),
			},
			{
				"Console",
				r.TypeOf((*string)(nil)).Elem(),
			},
			{
				"Status",
				r.TypeOf((*string)(nil)).Elem(),
			},
		},
		PeerPackage: "github.com/fpawel/anbus/internal/peer",
	})

}
