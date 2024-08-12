package schema

type Segment struct {
	s *Schema

	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Sample      string          `json:"sample"`
	Chapters    []string        `json:"chapters"`
	Fields      []*SegmentField `json:"fields"`
}

func (s *Segment) GetChapters() []*Chapter {
	var cs []*Chapter
	for _, c := range s.Chapters {
		cs = append(cs, s.s.Chapter(c))
	}

	return cs
}

func (s *Segment) GetFields() []*SegmentField {
	var fs []*SegmentField
	for _, f := range s.Fields {
		nf := f
		nf.s = s.s
		fs = append(fs, nf)
	}

	return fs
}

func (s *Segment) Field(id string) *SegmentField {
	for _, f := range s.Fields {
		if f.ID == id {
			f.s = s.s

			return f
		}
	}

	return nil
}

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

func (f *SegmentField) Required() bool {
	return f.MinRepetitions > 0
}

func (f *SegmentField) Repeatable() bool {
	return f.MaxRepetitions != 1
}

func (f *SegmentField) Table() *Table {
	if f.s == nil {
		return nil
	}

	return f.s.Table(f.TableID)
}

func (f *SegmentField) GetDataType() *DataType {
	if f.s == nil {
		return nil
	}

	return f.s.DataType(f.DataType)
}
