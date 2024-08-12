package schema

type schemaRegistry struct {
	schemas map[string]*Schema
}

var reg *schemaRegistry

func init() {
	reg = &schemaRegistry{
		schemas: make(map[string]*Schema),
	}
}

func RegisterSchema(s *Schema) {
	reg.schemas[s.Version()] = s
}

func Open(v string) *Schema {
	return reg.schemas[v]
}
