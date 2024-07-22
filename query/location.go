package query

import (
	"fmt"
	"strconv"
)

type Location struct {
	Segment      string
	SegmentRep   int
	Field        int
	FieldRep     int
	Component    int
	Subcomponent int
}

func ParseLocation(q string) (Location, error) {
	l := lex("query", q)

	loc := Location{}

	for i := l.nextItem(); i.typ != itemEOF; i = l.nextItem() {
		if i.typ == itemErr {
			return loc, fmt.Errorf("error parsing location: %s", i.val)
		}

		if i.typ == itemSeg {
			loc.Segment = i.val
			continue
		}

		iv, err := strconv.Atoi(i.val)
		if err != nil {
			return loc, err
		}

		switch i.typ {
		case itemSegIdx:
			loc.SegmentRep = iv
		case itemField:
			loc.Field = iv
		case itemRep:
			loc.FieldRep = iv
		case itemComp:
			loc.Component = iv
		case itemSub:
			loc.Subcomponent = iv
		}
	}

	return loc, nil
}
