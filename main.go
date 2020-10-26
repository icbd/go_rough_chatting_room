package main

import (
	"flag"
	"github.com/icbd/go_rough_chatting_room/client"
	"github.com/icbd/go_rough_chatting_room/server"
)

var addr = flag.String("addr", "localhost:8000", "Host and Port")
var serverMode = flag.Bool("server", false, "Default is ClientMode. Set this tag to turn on ServerMode")

func main() {
	flag.Parse()

	if *serverMode {
		server.New(*addr)
	} else {
		client.New(*addr)
	}
}
