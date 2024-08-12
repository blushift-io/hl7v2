package schema

type Message struct {
	s           *Schema
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Sample      string            `json:"sample"`
	Chapters    []string          `json:"chapters"`
	Segments    []*MessageSegment `json:"segments"`
}

func (m *Message) GetChapters() []*Chapter {
	var cs []*Chapter
	for _, c := range m.Chapters {
		cs = append(cs, m.s.Chapter(c))
	}

	return cs
}

func (m *Message) Segment(id string) *MessageSegment {
	for _, s := range m.Segments {
		if s.ID == id {
			s.s = m.s

			return s
		}
	}

	return nil
}

type MessageSegment struct {
	s *Schema

	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Sequence       int               `json:"sequence"`
	MinRepetitions int               `json:"min_repetitions"`
	MaxRepetitions int               `json:"max_repetitions"`
	Group          bool              `json:"group"`
	Segments       []*MessageSegment `json:"segments"`
}

func (ms *MessageSegment) Required() bool {
	return ms.MinRepetitions > 0
}

func (ms *MessageSegment) Repeatable() bool {
	return ms.MaxRepetitions != 1
}

func (ms *MessageSegment) Segment() *Segment {
	if ms.s == nil {
		return nil
	}

	if ms.Group {
		return nil
	}

	return ms.s.Segment(ms.ID)
}
