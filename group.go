package hl7v2

import "github.com/blushift-io/hl7v2/query"

// SegmentGroup represents a logical collection of HL7 segments.
type SegmentGroup struct {
	name     string
	parent   Element
	segments []*Segment
	index    []string
	count    map[string]int
}

func newSegmentGroup(name string, parent Element, segments []*Segment) *SegmentGroup {
	segGroup := &SegmentGroup{
		name:     name,
		parent:   parent,
		segments: segments,
		index:    make([]string, len(segments)),
		count:    make(map[string]int),
	}

	for i, seg := range segments {
		segGroup.index[i] = seg.Name()
		segGroup.count[seg.Name()]++
	}

	return segGroup
}

// Type returns the element type for SegmentGroup.
func (g *SegmentGroup) Type() ElementType {
	return ElementSegmentGroup
}

// Name returns the group name.
func (g *SegmentGroup) Name() string {
	return g.name
}

// Delimiters returns the message delimiters.
func (g *SegmentGroup) Delimiters() *Delimiters {
	return g.parent.Delimiters()
}

// Parent returns the parent element.
func (g *SegmentGroup) Parent() Element {
	return g.parent
}

// Children returns the segments in the group as child elements.
func (g *SegmentGroup) Children() []Element {
	return makeElements(g.segments...)
}

// Length returns the number of segments in the group.
func (g *SegmentGroup) Length() int {
	return len(g.segments)
}

// Position returns the position of the segment group.
func (g *SegmentGroup) Position() int {
	return 0
}

// Location returns the query Location of the segment group.
func (g *SegmentGroup) Location() query.Location {
	return g.parent.Location()
}

// GetLocation resolves an element inside the segment group by location.
func (g *SegmentGroup) GetLocation(_ query.Location) (Element, error) {
	panic("not implemented") // TODO: Implement
}

// Value returns the combined Value of all segments in the group.
func (g *SegmentGroup) Value(escape ...bool) Value {
	var b [][]byte

	for _, seg := range g.segments {
		b = append(b, seg.Value().Bytes())
	}

	return NewValue(g.parent.Delimiters().Join(b, SegmentDelimiter))
}
