package schema

// Segments represents a slice of Segment pointers.
type Segments []*Segment

// Len returns the number of segments in the collection.
func (ss Segments) Len() int {
	return len(ss)
}

// Less reports whether the segment at index i sorts before the segment at index j by ID.
func (ss Segments) Less(i, j int) bool {
	return ss[i].ID < ss[j].ID
}

// Swap exchanges the elements at indices i and j.
func (ss Segments) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

// Segment represents an HL7 segment definition.
type Segment struct {
	s *Schema

	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Sample      string          `json:"sample"`
	Chapters    []string        `json:"chapters"`
	Fields      []*SegmentField `json:"fields"`
}

// Version returns the schema version for the segment.
func (s *Segment) Version() string {
	if s.s == nil {
		return ""
	}

	return s.s.Version()
}

// GetChapters returns the chapters associated with the segment.
func (s *Segment) GetChapters() []*Chapter {
	var cs []*Chapter
	for _, c := range s.Chapters {
		cs = append(cs, s.s.Chapter(c))
	}

	return cs
}

// GetFields returns all field definitions belonging to the segment.
func (s *Segment) GetFields() []*SegmentField {
	var fs []*SegmentField
	for _, f := range s.Fields {
		nf := f
		nf.s = s.s
		fs = append(fs, nf)
	}

	return fs
}

// Field retrieves a segment field by its field ID.
func (s *Segment) Field(id string) *SegmentField {
	for _, f := range s.Fields {
		if f.ID == id {
			f.s = s.s

			return f
		}
	}

	return nil
}

// SegmentField represents a field entry within an HL7 segment definition.
type SegmentField struct {
	s *Schema

	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	Position       string `json:"position"`
	Length         int    `json:"length"`
	DataType       string `json:"dataType"`
	DataTypeName   string `json:"dataTypeName"`
	MinRepetitions int    `json:"min_repetitions"`
	MaxRepetitions int    `json:"max_repetitions"`
	TableID        string `json:"tableId"`
	TableName      string `json:"tableName"`
}

// Required reports whether the segment field is mandatory.
func (f *SegmentField) Required() bool {
	return f.MinRepetitions > 0
}

// Repeatable reports whether the segment field can repeat.
func (f *SegmentField) Repeatable() bool {
	return f.MaxRepetitions != 1
}

// Table returns the associated Table schema object for the field, if available.
func (f *SegmentField) Table() *Table {
	if f.s == nil {
		return nil
	}

	return f.s.Table(f.TableID)
}

// GetDataType returns the DataType schema object associated with the field.
func (f *SegmentField) GetDataType() *DataType {
	if f.s == nil {
		return nil
	}

	return f.s.DataType(f.DataType)
}
