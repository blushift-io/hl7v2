package query

import (
	"fmt"
	"regexp"
	"strings"
)

//TODO: Add support for named groups

const (
	optBegin = '['
	optEnd   = ']'
	rptBegin = '{'
	rptEnd   = '}'
	grpBegin = '('
	grpEnd   = ')'
)

var (
	matchSegment = regexp.MustCompile("^([A-Z])([A-Z])([A-Z]|[0-9])$")
)

func ParseGrammars(q string) (Grammars, error) {
	if len(q) == 0 {
		return nil, fmt.Errorf("grammar query cannot be empty")
	}

	return parseGrammar(q)
}

type Grammar struct {
	ID        string
	Optional  bool
	Repeating bool
	Group     bool
	Children  []Grammar
}

func (g Grammar) Validate(msg []string) error {
	idx, ok := matchInSet(g.ID, msg)
	if !g.Optional && !ok {
		return fmt.Errorf("missing required segment %s", g.ID)
	}

	for _, ch := range g.Children {
		if err := ch.Validate(msg[idx:]); err != nil {
			return err
		}
	}

	return nil
}

func (g Grammar) Select(msg []string) ([]string, error) {
	var res []string

	idx, ok := matchInSet(g.ID, msg)
	if !g.Optional && !ok {
		return nil, fmt.Errorf("missing required segment %s", g.ID)
	}

	if ok {
		res = append(res, g.ID)
	}

	for _, ch := range g.Children {
		segs, err := ch.Select(msg[idx:])
		if err != nil {
			return nil, err
		}

		res = append(res, segs...)
	}

	return res, nil
}

func (g Grammar) String() string {
	var b strings.Builder
	if g.Group {
		b.WriteString(g.ID)
	}

	if g.Optional {
		b.WriteByte(optBegin)
	}

	if g.Repeating {
		b.WriteByte(rptBegin)
	}

	if g.Group {
		j := make([]string, len(g.Children))
		for i, ch := range g.Children {
			j[i] = ch.String()
		}

		b.WriteString(strings.Join(j, " "))
	} else {
		b.WriteString(g.ID)
	}

	if g.Repeating {
		b.WriteByte(rptEnd)
	}

	if g.Optional {
		b.WriteByte(optEnd)
	}

	return b.String()
}

type Grammars []Grammar

func (g Grammars) Validate(msg []string) error {
	for _, m := range g {
		if err := m.Validate(msg); err != nil {
			return err
		}
	}

	return nil
}

func (g Grammars) Select(msg []string) ([]string, error) {
	var res []string
	for _, m := range g {
		r, err := m.Select(msg)
		if err != nil {
			return nil, err
		}

		res = append(res, r...)
	}

	return res, nil
}

func (g Grammars) String() string {
	j := make([]string, len(g))
	for i, m := range g {
		j[i] = m.String()
	}

	return strings.Join(j, " ")
}

func parseGrammar(expr string) ([]Grammar, error) {
	idx := 0
	var res []Grammar

	openers := []byte{optBegin, rptBegin, grpBegin}
	closers := []byte{optEnd, rptEnd, grpEnd}

	for idx < len(expr) {
		ch := expr[idx]

		switch {
		case ch == ' ':
			idx++
			continue
		case setContains(ch, openers):
			y, _ := matchInSet(ch, openers)
			open, close := openers[y], closers[y]
			adv, iExpr := parseInnerGrammar(expr, idx, open, close)
			idx = adv
			isOpt := open == optBegin
			isRep := open == rptBegin
			//isGrp := open == grpBegin

			if isOpt && iExpr[0] == rptBegin {
				isRep = true
				iExpr = iExpr[1 : len(iExpr)-1]
			}

			if len(iExpr) == 3 {
				res = append(res, Grammar{
					ID:        iExpr,
					Optional:  isOpt,
					Repeating: isRep,
				})

				continue
			}

			inner, err := parseGrammar(iExpr)
			if err != nil {
				return nil, err
			}

			if len(inner) == 0 {
				continue
			}

			gr := Grammar{
				ID:        inner[0].ID,
				Optional:  isOpt,
				Repeating: isRep,
			}

			if len(inner) > 1 {
				gr.Children = inner[1:]
			}

			res = append(res, gr)
		default:
			adv, id := parserSegmentID(expr, idx)
			if adv == -1 {
				return nil, fmt.Errorf("error parsing segment")
			}

			idx = adv
			res = append(res, Grammar{
				ID: id,
			})
		}
	}

	return res, nil
}

func parserSegmentID(expr string, idx int) (int, string) {
	if len(expr) < idx+3 {
		return -1, ""
	}

	seg := expr[idx : idx+3]
	if matchSegment.MatchString(seg) {
		return idx + 3, seg
	}

	return -1, ""
}

func parseInnerGrammar(expr string, idx int, open, close byte) (int, string) {
	depth := 0
	b := &strings.Builder{}

	idx++
	for idx < len(expr) {
		ch := expr[idx]
		if depth == 0 && ch == close {
			break
		}

		if ch == open {
			depth++
		}

		if ch == close {
			depth--
		}

		b.WriteByte(ch)
		idx++
	}

	idx++

	return idx, b.String()
}

func setContains[T comparable](v T, set []T) bool {
	_, ok := matchInSet(v, set)
	return ok
}

func matchInSet[T comparable](v T, set []T) (int, bool) {
	for i, t := range set {
		if v == t {
			return i, true
		}
	}

	return -1, false
}
