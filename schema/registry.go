package schema

import (
	"errors"
)

var (
	// ErrSchemaNotFound indicates that a requested HL7 schema version was not found in the registry.
	ErrSchemaNotFound = errors.New("schema not found")
)

type schemaRegistry struct {
	schemas map[string]*Schema
}

var reg *schemaRegistry

func init() {
	reg = &schemaRegistry{
		schemas: make(map[string]*Schema),
	}
}

// All returns all registered HL7 schemas.
func All() []*Schema {
	s := make([]*Schema, 0, len(reg.schemas))
	for _, v := range reg.schemas {
		s = append(s, v)
	}

	return s
}

// Register registers an HL7 schema in the global registry using its version string.
func Register(s *Schema) {
	reg.schemas[s.Version()] = s
}

// Open returns the registered HL7 schema for the specified version, or nil if not found.
func Open(v string) *Schema {
	return reg.schemas[v]
}
