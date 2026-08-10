// Package main provides a code generation tool for HL7v2 schema Go code.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/versions/generator"
	"github.com/urfave/cli/v3"

	_ "github.com/blushift-io/hl7v2/schema/spec/v21"
	_ "github.com/blushift-io/hl7v2/schema/spec/v22"
	_ "github.com/blushift-io/hl7v2/schema/spec/v23"
	_ "github.com/blushift-io/hl7v2/schema/spec/v231"
	_ "github.com/blushift-io/hl7v2/schema/spec/v24"
	_ "github.com/blushift-io/hl7v2/schema/spec/v25"
	_ "github.com/blushift-io/hl7v2/schema/spec/v251"
	_ "github.com/blushift-io/hl7v2/schema/spec/v26"
	_ "github.com/blushift-io/hl7v2/schema/spec/v27"
	_ "github.com/blushift-io/hl7v2/schema/spec/v271"
	_ "github.com/blushift-io/hl7v2/schema/spec/v28"
)

const outputBase = "./versions"

func main() {
	cmd := &cli.Command{
		Name: "gen_hl7v2",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Value: "all",
				Usage: "HL7v2 version to generate (e.g. 2.5.1 or 'all')",
			},
		},
		Action: genAction,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func genAction(ctx context.Context, cmd *cli.Command) error {
	vf := strings.TrimSpace(cmd.String("version"))
	if vf == "" || strings.EqualFold(vf, "all") {
		return generateAllVersions()
	}

	return generateVersion(vf)
}

func generateVersion(vStr string) error {
	gen, err := generator.New(vStr)
	if err != nil {
		return err
	}

	log.Printf("Generating HL7v2 schema for version %s...", vStr)
	if err := gen.GenerateAll(outputBase); err != nil {
		return err
	}
	log.Printf("Successfully generated code for version %s.", vStr)

	return nil
}

func generateAllVersions() error {
	for _, v := range hl7v2.VersionValues() {
		if v == hl7v2.VersionUnknown {
			continue
		}
		vStr := v.String()
		gen, err := generator.New(vStr)
		if err != nil {
			log.Printf("Skipping version %s: %v", vStr, err)
			continue
		}
		log.Printf("Generating HL7v2 schema for version %s...", vStr)
		if err := gen.GenerateAll(outputBase); err != nil {
			return fmt.Errorf("failed generating version %s: %w", vStr, err)
		}
	}
	log.Println("Finished generating code for all versions.")
	return nil
}
