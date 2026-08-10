package schema

// TableType represents the category or source of an HL7 table.
type TableType int

const (
	// TableTypeUnknown represents an unknown or unspecified table type.
	TableTypeUnknown TableType = iota //Unknown
	// TableTypeUser represents a user-defined table.
	TableTypeUser //User
	// TableTypeHL7 represents an HL7 standard table.
	TableTypeHL7 //HL7
	// TableTypeLocal represents a local system table.
	TableTypeLocal //Local
	// TableTypePreLoaded represents a pre-loaded reference table.
	TableTypePreLoaded //PreLoaded
)

// Tables represents a slice of Table pointers.
type Tables []*Table

// Len returns the number of tables in the collection.
func (ts Tables) Len() int {
	return len(ts)
}

// Less reports whether the table at index i sorts before the table at index j by ID.
func (ts Tables) Less(i, j int) bool {
	return ts[i].ID < ts[j].ID
}

// Swap exchanges the elements at indices i and j.
func (ts Tables) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

// Table represents an HL7 table definition containing value mappings.
type Table struct {
	s *Schema

	ID       string       `json:"id"`
	Type     TableType    `json:"type"`
	Name     string       `json:"name"`
	Chapters []string     `json:"chapters"`
	Entries  []TableEntry `json:"entries"`
}

// GetChapters returns the chapters associated with the table.
func (t *Table) GetChapters() []*Chapter {
	var cs []*Chapter
	for _, c := range t.Chapters {
		cs = append(cs, t.s.Chapter(c))
	}

	return cs
}

// Entry retrieves a table entry matching the provided value string.
func (t *Table) Entry(value string) (*TableEntry, bool) {
	for _, e := range t.Entries {
		if e.Value == value {
			return &e, true
		}
	}

	return nil, false
}

// TableEntry represents a single code/value mapping within an HL7 table.
type TableEntry struct {
	Value       string `json:"value"`
	Description string `json:"description"`
	Comment     string `json:"comment,omitempty"`
}
