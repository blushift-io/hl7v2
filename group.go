package hl7v2

import "github.com/blushift-io/hl7v2/query"

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

func (g *SegmentGroup) Type() ElementType {
	return ElementSegmentGroup
}

func (g *SegmentGroup) Name() string {
	return g.name
}

func (g *SegmentGroup) Delimiters() *Delimiters {
	return g.parent.Delimiters()
}

func (g *SegmentGroup) Parent() Element {
	return g.parent
}

func (g *SegmentGroup) Children() []Element {
	return makeElements(g.segments...)
}

func (g *SegmentGroup) Length() int {
	return len(g.segments)
}

func (g *SegmentGroup) Position() int {
	return 0
}

func (g *SegmentGroup) Location() query.Location {
	return g.parent.Location()
}

func (g *SegmentGroup) GetLocation(_ query.Location) (Element, error) {
	panic("not implemented") // TODO: Implement
}

func (g *SegmentGroup) Value() Value {
	var b [][]byte

	for _, seg := range g.segments {
		b = append(b, seg.Value().Bytes())
	}

	return NewValue(g.parent.Delimiters().Join(b, SegmentDelimiter))
}
