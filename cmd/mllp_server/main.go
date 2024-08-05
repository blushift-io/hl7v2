package main

import (
	"log"

	"github.com/blushift-io/hl7v2/net/tcp"
)

func main() {
	h := &handler{
		savePath: "./tmp",
	}

	s, err := tcp.NewServer(h, tcp.WithAddresses(":2525"))
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
