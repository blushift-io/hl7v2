package schema

type Chapters []*Chapter

func (c Chapters) Len() int {
	return len(c)
}

func (c Chapters) Less(i int, j int) bool {
	return c[i].ID < c[j].ID
}

func (c Chapters) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

type Chapter struct {
	s *Schema

	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

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
