package schema

import (
	"errors"
)

var (
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

func All() []*Schema {
	s := make([]*Schema, 0, len(reg.schemas))
	for _, v := range reg.schemas {
		s = append(s, v)
	}

	return s
}

func Register(s *Schema) {
	reg.schemas[s.Version()] = s
}

func Open(v string) *Schema {
	return reg.schemas[v]
}
