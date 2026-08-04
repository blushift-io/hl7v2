package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/blushift-io/hl7v2/net/tcp"
	"github.com/oklog/run"
)

func main() {
	h := &handler{
		savePath: "./tmp",
	}

	s, err := tcp.NewServer(h, tcp.WithAddresses(":2525"))
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	g := &run.Group{}
	ctx := context.Background()

	g.Add(s.Start, func(e error) {
		if err := s.Shutdown(); err != nil {
			log.Fatal(err)
		}
	})

	g.Add(run.SignalHandler(ctx, os.Interrupt, os.Kill))

	if err := g.Run(); err != nil {
		if !errors.Is(err, run.ErrSignal) {
			log.Fatal(err)
		}
	}

	log.Println("exited")
}
