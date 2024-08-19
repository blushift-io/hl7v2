package schema

type TableType int

const (
	TableTypeUnknown   TableType = iota //Unknown
	TableTypeUser                       //User
	TableTypeHL7                        //HL7
	TableTypeLocal                      //Local
	TableTypePreLoaded                  //PreLoaded
)

type Tables []*Table

func (ts Tables) Len() int {
	return len(ts)
}

func (ts Tables) Less(i, j int) bool {
	return ts[i].ID < ts[j].ID
}

func (ts Tables) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

type Table struct {
	s *Schema

	ID       string       `json:"id"`
	Type     TableType    `json:"type"`
	Name     string       `json:"name"`
	Chapters []string     `json:"chapters"`
	Entries  []TableEntry `json:"entries"`
}

func (t *Table) GetChapters() []*Chapter {
	var cs []*Chapter
	for _, c := range t.Chapters {
		cs = append(cs, t.s.Chapter(c))
	}

	return cs
}

func (t *Table) Entry(value string) (*TableEntry, bool) {
	for _, e := range t.Entries {
		if e.Value == value {
			return &e, true
		}
	}

	return nil, false
}

type TableEntry struct {
	Value       string `json:"value"`
	Description string `json:"description"`
	Comment     string `json:"comment,omitempty"`
}
