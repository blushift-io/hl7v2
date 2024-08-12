package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Schema struct {
	version   string
	chapters  map[string]*Chapter
	tables    map[string]*Table
	dataTypes map[string]*DataType
	messages  map[string]*Message
	segments  map[string]*Segment
}

func LoadJSONSchema(p string) (*Schema, error) {
	sf := filepath.Join(p, "schema.json")
	sm := make(map[string]any)
	if err := loadJSON(sf, &sm); err != nil {
		return nil, err
	}

	ver, ok := sm["version"].(string)
	if !ok {
		return nil, fmt.Errorf("missing version in schema")
	}

	s := &Schema{
		version:   ver,
		chapters:  make(map[string]*Chapter),
		tables:    make(map[string]*Table),
		dataTypes: make(map[string]*DataType),
		messages:  make(map[string]*Message),
		segments:  make(map[string]*Segment),
	}

	cf := filepath.Join(p, "chapters.json")
	var ch []*Chapter
	if err := loadJSON(cf, &ch); err != nil {
		return nil, err
	}

	for _, c := range ch {
		c.s = s
		s.chapters[c.ID] = c
	}

	tf := filepath.Join(p, "tables.json")
	var ta []*Table
	if err := loadJSON(tf, &ta); err != nil {
		return nil, err
	}

	for _, t := range ta {
		t.s = s
		s.tables[t.ID] = t
	}

	dtf := filepath.Join(p, "datatypes.json")
	var dt []*DataType
	if err := loadJSON(dtf, &dt); err != nil {
		return nil, err
	}

	for _, d := range dt {
		d.s = s
		s.dataTypes[d.ID] = d
	}

	mf := filepath.Join(p, "messages.json")
	var ms []*Message
	if err := loadJSON(mf, &ms); err != nil {
		return nil, err
	}

	for _, m := range ms {
		m.s = s
		s.messages[m.ID] = m
	}

	sf = filepath.Join(p, "segments.json")
	var sg []*Segment
	if err := loadJSON(sf, &sg); err != nil {
		return nil, err
	}

	for _, seg := range sg {
		seg.s = s
		s.segments[seg.ID] = seg
	}

	return s, nil
}

func loadJSON(p string, v any) error {
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (s *Schema) Version() string {
	return s.version
}

func (s *Schema) Chapters() []*Chapter {
	var cs []*Chapter
	for _, c := range s.chapters {
		cs = append(cs, c)
	}

	return cs
}

func (s *Schema) Chapter(id string) *Chapter {
	return s.chapters[id]
}

func (s *Schema) Tables() []*Table {
	var ts []*Table
	for _, t := range s.tables {
		ts = append(ts, t)
	}

	return ts
}

func (s *Schema) Table(id string) *Table {
	return s.tables[id]
}

func (s *Schema) DataTypes() []*DataType {
	var ds []*DataType
	for _, d := range s.dataTypes {
		ds = append(ds, d)
	}

	return ds
}

func (s *Schema) DataType(id string) *DataType {
	return s.dataTypes[id]
}

func (s *Schema) Messages() []*Message {
	var ms []*Message
	for _, m := range s.messages {
		ms = append(ms, m)
	}

	return ms
}

func (s *Schema) Message(id string) *Message {
	return s.messages[id]
}

func (s *Schema) Segments() []*Segment {
	var ss []*Segment
	for _, s := range s.segments {
		ss = append(ss, s)
	}

	return ss
}

func (s *Schema) Segment(id string) *Segment {
	return s.segments[id]
}
