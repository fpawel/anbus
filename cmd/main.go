package main

import (
	"flag"
	"github.com/fpawel/anbus/internal/work"
)

func main() {
	ensurePeer := true
	flag.BoolVar(&ensurePeer, "ensure-peer", true, "ensure peer application")
	flag.Parse()
	work.Main()
}
