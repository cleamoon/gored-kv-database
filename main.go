package main

import (
	"log"
	"net"

	"github.com/cleamoon/gored-kv-database/internal/server"
	"github.com/cleamoon/gored-kv-database/pkg/kvstore"
)

func main() {
	line, err := net.Listen("tcp", ":6789")
	if err != nil {
		log.Fatal(err)
	}
	defer line.Close()

	kv := kvstore.NewKVStore()

	for {
		conn, err := line.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.HandleConnection(conn, kv)
	}
}