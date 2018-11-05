package main

import (
	"flag"
	"github.com/fpawel/anbus/internal/work"
)

func main() {
	mustRunPeer := true
	flag.BoolVar(&mustRunPeer, "must-run-peer", true, "ensure peer application")
	flag.Parse()
	work.Main(mustRunPeer)
}
