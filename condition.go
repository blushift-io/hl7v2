package hl7v2

import (
	"strings"

	"github.com/blushift-io/hl7v2/query"
)

// ConditionFunc is a predicate function that evaluates an HL7 Element.
type ConditionFunc func(Element) bool

// Condition represents a named rule for testing HL7 elements.
type Condition struct {
	Name  string
	Check ConditionFunc
}

// Has returns a Condition that checks if the specified query location exists and is non-empty.
func Has(q string) (*Condition, error) {
	loc, err := query.ParseLocation(q)
	if err != nil {
		return nil, err
	}

	return &Condition{
		Name: "Has",
		Check: func(e Element) bool {
			el, err := e.GetLocation(loc)
			if err != nil {
				return false
			}

			if el == nil {
				return false
			}

			return !el.Value().Empty()
		},
	}, nil
}

// ValueContains returns a Condition that checks if the value at the location contains the substring v.
func ValueContains(q, v string) (*Condition, error) {
	loc, err := query.ParseLocation(q)
	if err != nil {
		return nil, err
	}

	return &Condition{
		Name: "ValueContains",
		Check: func(e Element) bool {
			el, err := e.GetLocation(loc)
			if err != nil {
				return false
			}

			if el == nil {
				return false
			}

			return strings.Contains(el.Value().String(), v)
		},
	}, nil
}

// ValueContainsAny returns a Condition that checks if the value at the location contains any of the substrings in v.
func ValueContainsAny(q string, v ...string) (*Condition, error) {
	loc, err := query.ParseLocation(q)
	if err != nil {
		return nil, err
	}

	return &Condition{
		Name: "ValueContainsAny",
		Check: func(e Element) bool {
			el, err := e.GetLocation(loc)
			if err != nil {
				return false
			}

			if el == nil {
				return false
			}

			for _, s := range v {
				if strings.Contains(el.Value().String(), s) {
					return true
				}
			}

			return false
		},
	}, nil
}
