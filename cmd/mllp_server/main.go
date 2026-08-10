// Package main implements a server CLI application that receives, processes, saves, and forwards HL7v2 messages via MLLP.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/blushift-io/hl7v2/net/tcp"
	"github.com/oklog/run"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "mllp_server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				Aliases:   []string{"c"},
				Usage:     "config file",
				Value:     "config.yaml",
				TakesFile: true,
			},
		},
		Action: runServer,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx context.Context, cmd *cli.Command) error {
	conf, err := loadConfig(cmd)
	if err != nil {
		if conf == nil {
			return fmt.Errorf("unable to load config: %w", err)
		}

		log.Printf("error loading config, using default config: %v", err)
	}

	h, err := newHandler(conf.Actions)
	if err != nil {
		return fmt.Errorf("unable to create handler: %w", err)
	}

	addr := fmt.Sprintf(":%d", conf.Port)
	s, err := tcp.NewServer(h, tcp.WithAddresses(addr))
	if err != nil {
		return fmt.Errorf("unable to create server: %w", err)
	}

	g := &run.Group{}

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

	return nil
}

func loadConfig(cmd *cli.Command) (*config, error) {
	fileName := cmd.String("config")
	return readConfig(fileName)
}
