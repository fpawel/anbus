package app

import (
	"fmt"
	"github.com/fpawel/anbus/internal/api"
	"github.com/fpawel/dseries"
	"github.com/fpawel/gohelp/must"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"golang.org/x/sys/windows/registry"
	"net"
	"net/http"
	"net/rpc"
)

func startHttpServer() func() {

	for _, svcObj := range []interface{}{
		&api.ConfigSvc{chConfigChanged},
		new(dseries.ChartsSvc),
		&api.TaskSvc{taskRunner{}},
	} {
		must.AbortIf(rpc.Register(svcObj))
	}

	// Server provide a HTTP transport on /rpc endpoint.
	http.Handle("/rpc", jsonrpc2.HTTPHandler(nil))

	http.HandleFunc("/chart", dseries.HandleRequestChart)

	http.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("assets"))))

	srv := new(http.Server)
	lnHTTP, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	addr := "http://" + lnHTTP.Addr().String()
	fmt.Printf("%s/report?party_id=last\n", addr)
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `anbus\http`, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	if err := key.SetStringValue("addr", addr); err != nil {
		panic(err)
	}
	log.ErrIfFail(key.Close)

	go func() {
		err := srv.Serve(lnHTTP)
		if err == http.ErrServerClosed {
			return
		}
		log.PrintErr(err)
		log.ErrIfFail(lnHTTP.Close)
	}()

	return func() {
		if err := srv.Shutdown(ctxApp); err != nil {
			log.PrintErr(err)
		}
	}
}
