package schema

import (
	"encoding/json"
	"sort"
	"sync"
)

//TODO: Complete GoDoc to explain schema linking via functional getter vs struct field access

//go:generate enumer -type=TableType,DataTypeType -json -text -sql -yaml -linecomment -output=schema_enums.go

type Schema struct {
	l sync.RWMutex

	version   string
	chapters  map[string]*Chapter
	tables    map[string]*Table
	dataTypes map[string]*DataType
	messages  map[string]*Message
	segments  map[string]*Segment
}

func NewSchema(v string) *Schema {
	return &Schema{
		version:   v,
		chapters:  make(map[string]*Chapter),
		tables:    make(map[string]*Table),
		dataTypes: make(map[string]*DataType),
		messages:  make(map[string]*Message),
		segments:  make(map[string]*Segment),
	}
}

func (s *Schema) AddChapters(ch ...*Chapter) {
	s.l.Lock()
	defer s.l.Unlock()

	for _, c := range ch {
		c.s = s
		s.chapters[c.ID] = c
	}
}

func (s *Schema) AddTables(t ...*Table) {
	s.l.Lock()
	defer s.l.Unlock()

	for _, ta := range t {
		ta.s = s
		s.tables[ta.ID] = ta
	}
}

func (s *Schema) AddDataTypes(d ...*DataType) {
	s.l.Lock()
	defer s.l.Unlock()

	for _, dt := range d {
		dt.s = s
		s.dataTypes[dt.ID] = dt
	}
}

func (s *Schema) AddMessages(m ...*Message) {
	s.l.Lock()
	defer s.l.Unlock()

	for _, ms := range m {
		ms.s = s
		s.messages[ms.ID] = ms
	}
}

func (s *Schema) AddSegments(se ...*Segment) {
	s.l.Lock()
	defer s.l.Unlock()

	for _, seg := range se {
		seg.s = s
		s.segments[seg.ID] = seg
	}
}

func (s *Schema) Version() string {
	return s.version
}

func (s *Schema) Chapters() []*Chapter {
	s.l.RLock()
	defer s.l.RUnlock()

	var cs Chapters
	for _, c := range s.chapters {
		cs = append(cs, c)
	}

	sort.Sort(cs)

	return cs
}

func (s *Schema) Chapter(id string) *Chapter {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.chapters[id]
}

func (s *Schema) Tables() []*Table {
	s.l.RLock()
	defer s.l.RUnlock()

	var ts Tables
	for _, t := range s.tables {
		ts = append(ts, t)
	}

	sort.Sort(ts)

	return ts
}

func (s *Schema) Table(id string) *Table {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.tables[id]
}

func (s *Schema) DataTypes() []*DataType {
	s.l.RLock()
	defer s.l.RUnlock()

	var ds DataTypes
	for _, d := range s.dataTypes {
		ds = append(ds, d)
	}

	sort.Sort(ds)

	return ds
}

func (s *Schema) DataType(id string) *DataType {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.dataTypes[id]
}

func (s *Schema) Messages() []*Message {
	s.l.RLock()
	defer s.l.RUnlock()

	var ms Messages
	for _, m := range s.messages {
		ms = append(ms, m)
	}

	sort.Sort(ms)
	return ms
}

func (s *Schema) Message(id string) *Message {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.messages[id]
}

func (s *Schema) Segments() []*Segment {
	s.l.RLock()
	defer s.l.RUnlock()

	var ss Segments
	for _, s := range s.segments {
		ss = append(ss, s)
	}

	sort.Sort(ss)
	return ss
}

func (s *Schema) Segment(id string) *Segment {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.segments[id]
}

type jsonSchema struct {
	Version   string      `json:"version"`
	Chapters  []*Chapter  `json:"chapters"`
	Tables    []*Table    `json:"tables"`
	DataTypes []*DataType `json:"datatypes"`
	Messages  []*Message  `json:"messages"`
	Segments  []*Segment  `json:"segments"`
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonSchema{
		Version:   s.version,
		Chapters:  s.Chapters(),
		Tables:    s.Tables(),
		DataTypes: s.DataTypes(),
		Messages:  s.Messages(),
		Segments:  s.Segments(),
	})
}
