package v26

import (
	_ "embed"

	"github.com/blushift-io/hl7v2/schema"
)

//go:embed schema.json
var version []byte

//go:embed chapters.json
var chapters []byte

//go:embed tables.json
var tables []byte

//go:embed datatypes.json
var datatypes []byte

//go:embed messages.json
var messages []byte

//go:embed segments.json
var segments []byte

func init() {
	if err := LoadSchema(); err != nil {
		panic(err)
	}
}

func LoadSchema() error {
	spec := &schema.JSONSpec{
		Schema:    version,
		Chapters:  chapters,
		Tables:    tables,
		DataTypes: datatypes,
		Messages:  messages,
		Segments:  segments,
	}

	sch, err := spec.Build()
	if err != nil {
		return err
	}

	schema.Register(sch)

	return nil
}
