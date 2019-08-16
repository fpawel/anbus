package main

import (
	"flag"
	"github.com/fpawel/anbus/internal/app"
)

func main() {
	mustRunPeer := true
	flag.BoolVar(&mustRunPeer, "must-run-peer", true, "ensure peer application")
	flag.Parse()
	app.Main(mustRunPeer)
}
