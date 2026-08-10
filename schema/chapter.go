package schema

// Chapters represents a slice of Chapter pointers.
type Chapters []*Chapter

// Len returns the number of chapters in the collection.
func (c Chapters) Len() int {
	return len(c)
}

// Less reports whether the chapter at index i sorts before the chapter at index j by ID.
func (c Chapters) Less(i int, j int) bool {
	return c[i].ID < c[j].ID
}

// Swap exchanges the elements at indices i and j.
func (c Chapters) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

// Chapter represents an HL7 specification chapter.
type Chapter struct {
	s *Schema

	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Messages returns all messages associated with this chapter in the schema.
func (c *Chapter) Messages() []*Message {
	if c.s == nil {
		return nil
	}

	var ms []*Message
	for _, m := range c.s.Messages() {
		for _, ch := range m.Chapters {
			if ch == c.ID {
				ms = append(ms, m)
				break
			}
		}
	}

	return ms
}
