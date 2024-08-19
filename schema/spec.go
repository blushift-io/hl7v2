package schema

import (
	"encoding/json"
	"fmt"
)

type JSONSpec struct {
	Schema    []byte
	Chapters  []byte
	Tables    []byte
	DataTypes []byte
	Messages  []byte
	Segments  []byte
}

func (spec *JSONSpec) MustBuild() *Schema {
	s, err := spec.Build()
	if err != nil {
		panic(err)
	}

	return s
}

func (spec *JSONSpec) Build() (*Schema, error) {
	sm := make(map[string]any)
	if err := json.Unmarshal(spec.Schema, &sm); err != nil {
		return nil, err
	}

	ver, ok := sm["version"].(string)
	if !ok {
		return nil, fmt.Errorf("missing version in schema")
	}

	s := NewSchema(ver)

	var ch []*Chapter
	if err := json.Unmarshal(spec.Chapters, &ch); err != nil {
		return nil, err
	}

	s.AddChapters(ch...)

	var ta []*Table
	if err := json.Unmarshal(spec.Tables, &ta); err != nil {
		return nil, err
	}

	s.AddTables(ta...)

	var dt []*DataType
	if err := json.Unmarshal(spec.DataTypes, &dt); err != nil {
		return nil, err
	}

	s.AddDataTypes(dt...)

	var ms []*Message
	if err := json.Unmarshal(spec.Messages, &ms); err != nil {
		return nil, err
	}

	s.AddMessages(ms...)

	var se []*Segment
	if err := json.Unmarshal(spec.Segments, &se); err != nil {
		return nil, err
	}

	s.AddSegments(se...)

	return s, nil
}
