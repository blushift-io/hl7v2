package schema

import "github.com/blushift-io/hl7v2/query"

type Messages []*Message

func (m Messages) Len() int {
	return len(m)
}

func (m Messages) Less(i int, j int) bool {
	return m[i].ID < m[j].ID
}

func (m Messages) Swap(i int, j int) {
	m[i], m[j] = m[j], m[i]
}

type Message struct {
	s           *Schema
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Sample      string            `json:"sample"`
	Chapters    []string          `json:"chapters"`
	Segments    []*MessageSegment `json:"segments"`
}

func (m *Message) Version() string {
	if m.s == nil {
		return ""
	}

	return m.s.version
}

func (m *Message) GetChapters() []*Chapter {
	var cs []*Chapter
	for _, c := range m.Chapters {
		cs = append(cs, m.s.Chapter(c))
	}

	return cs
}

func (m *Message) GetSegments() []*MessageSegment {
	var ss []*MessageSegment
	for _, s := range m.Segments {
		s.s = m.s
		ss = append(ss, s)
	}

	return ss
}

func (m *Message) Segment(id string) *MessageSegment {
	seg := matchSegment(id, m.Segments)
	if seg != nil {
		seg.s = m.s
	}

	return seg
}

func matchSegment(id string, segs []*MessageSegment) *MessageSegment {
	for _, s := range segs {
		if s.Group {
			if seg := matchSegment(id, s.Segments); seg != nil {
				return seg
			}
		}

		if s.ID == id || s.Name == id {
			return s
		}
	}

	return nil
}

func (m *Message) Grammar() query.Grammars {
	var g query.Grammars

	for _, s := range m.Segments {
		g = append(g, s.Grammar())
	}

	return g
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

func (ms *MessageSegment) Grammar() query.Grammar {
	var g query.Grammar

	if ms.Group {
		g = query.Grammar{
			ID:        ms.Name,
			Optional:  !ms.Required(),
			Repeating: ms.Repeatable(),
			Group:     ms.Group,
		}

		for _, s := range ms.Segments {
			g.Children = append(g.Children, s.Grammar())
		}

		return g
	}

	return query.Grammar{
		ID:        ms.ID,
		Optional:  !ms.Required(),
		Repeating: ms.Repeatable(),
		Group:     ms.Group,
	}
}
