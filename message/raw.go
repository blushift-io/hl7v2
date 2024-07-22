package message

import (
	"encoding/json"
	"fmt"

	"github.com/blushift-io/hl7v2"
)

type RawSubcomponent []byte
type RawComponent []RawSubcomponent
type RawRepetition []RawComponent
type RawField []RawRepetition
type RawSegment []RawField

type RawMessage struct {
	delims *hl7v2.Delimiters
	segs   []RawSegment
}

func (m *RawMessage) Segments() []RawSegment {
	return m.segs
}

func (m *RawMessage) Append(seg RawSegment) {
	m.segs = append(m.segs, seg)
}

func (m *RawMessage) JSON() ([]byte, error) {
	msg := jsonMessage{}

	for _, seg := range m.segs {
		msg.Segments = append(msg.Segments, extractJSONSegment(seg))
	}

	return json.MarshalIndent(msg, "", "  ")
}

type jsonMessage struct {
	Segments []jsonSegment `json:"segments"`
}

type jsonSegment struct {
	ID     string         `json:"id"`
	Fields map[string]any `json:"fields"`
}

func extractJSONSegment(seg RawSegment) jsonSegment {
	id := seg[0][0][0][0]
	flds := make(map[string]any)

	for fi, fld := range seg {
		if fi == 0 {
			continue
		}

		idx := fmt.Sprintf("%d", fi)
		for ri, rep := range fld {
			if len(fld) > 1 {
				idx = fmt.Sprintf("%s[%d]", idx, ri+1)
			}

			for ci, cmp := range rep {
				if len(rep) > 1 {
					idx = fmt.Sprintf("%s.%d", idx, ci+1)
				}

				for si, sub := range cmp {
					if len(cmp) > 1 {
						idx = fmt.Sprintf("%s-%d", idx, si+1)
					}

					flds[idx] = string(sub)
				}
			}
		}
	}

	return jsonSegment{
		ID:     string(id),
		Fields: flds,
	}
}
