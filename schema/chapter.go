package schema

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
