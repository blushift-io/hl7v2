package schema

import "github.com/blushift-io/hl7v2/query"

// Messages represents a slice of Message pointers.
type Messages []*Message

// Len returns the number of messages in the collection.
func (m Messages) Len() int {
	return len(m)
}

// Less reports whether the message at index i sorts before the message at index j by ID.
func (m Messages) Less(i int, j int) bool {
	return m[i].ID < m[j].ID
}

// Swap exchanges the elements at indices i and j.
func (m Messages) Swap(i int, j int) {
	m[i], m[j] = m[j], m[i]
}

// Message represents an HL7 v2 message structure definition.
type Message struct {
	s           *Schema
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Sample      string            `json:"sample"`
	Chapters    []string          `json:"chapters"`
	Segments    []*MessageSegment `json:"segments"`
}

// Version returns the HL7 v2 version of the message.
func (m *Message) Version() string {
	if m.s == nil {
		return ""
	}

	return m.s.version
}

// GetChapters returns the chapters associated with the message.
func (m *Message) GetChapters() []*Chapter {
	var cs []*Chapter
	for _, c := range m.Chapters {
		cs = append(cs, m.s.Chapter(c))
	}

	return cs
}

// GetSegments returns the list of segment definitions belonging to the message.
func (m *Message) GetSegments() []*MessageSegment {
	var ss []*MessageSegment
	for _, s := range m.Segments {
		s.s = m.s
		ss = append(ss, s)
	}

	return ss
}

// Segment retrieves a message segment or segment group by its ID or name.
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

// Grammar constructs and returns the query.Grammars rule collection for the message structure.
func (m *Message) Grammar() query.Grammars {
	var g query.Grammars

	for _, s := range m.Segments {
		g = append(g, s.Grammar())
	}

	return g
}

// MessageSegment represents a segment or segment group entry within an HL7 message schema.
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

// Required reports whether the segment or group is required in the message.
func (ms *MessageSegment) Required() bool {
	return ms.MinRepetitions > 0
}

// Repeatable reports whether the segment or group can repeat.
func (ms *MessageSegment) Repeatable() bool {
	return ms.MaxRepetitions != 1
}

// Segment returns the underlying Segment definition from the schema if this is not a group.
func (ms *MessageSegment) Segment() *Segment {
	if ms.s == nil {
		return nil
	}

	if ms.Group {
		return nil
	}

	return ms.s.Segment(ms.ID)
}

// Grammar constructs and returns a query.Grammar rule representing the message segment or group.
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
