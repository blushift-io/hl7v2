package hl7v2

import (
	"strings"

	"github.com/blushift-io/hl7v2/query"
)

type ConditionFunc func(Element) bool

type Condition struct {
	Name  string
	Check ConditionFunc
}

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
