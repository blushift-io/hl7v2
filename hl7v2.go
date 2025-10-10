package hl7v2

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structs"
)

//TODO: Escape/Unescape HL7 delimiters
//TODO: Implement MCF delayed acknowledgment
//TODO: Implement batch and file elements

type Marshaler interface {
	MarshalHL7() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalHL7(msg Queryable, v any) error
}

func Unmarshal(msg Queryable, v any) error {
	if m, ok := v.(Unmarshaler); ok {
		return m.UnmarshalHL7(msg, v)
	}

	if !structs.IsStruct(v) {
		return fmt.Errorf("hl7v2: Unmarshal expects a struct")
	}

	s := structs.New(v)

	for _, f := range s.Fields() {
		if !f.IsExported() {
			continue
		}

		q := f.Tag("hl7")
		if len(q) == 0 {
			continue
		}

		var ql []string
		if -1 < strings.Index(q, ",") {
			ql = strings.Split(q, ",")
		} else {
			ql = []string{q}
		}

		for _, loc := range ql {
			q := strings.Trim(loc, " ")

			val, err := msg.QueryValue(q)
			if err != nil {
				if errors.Is(err, ErrValueNotFound) || errors.Is(err, ErrElementNotFound) {
					continue
				}

				return err
			}

			if val.Empty() {
				continue
			}

			var setVal any
			switch f.Value().(type) {
			case string:
				setVal = val.String()
			case int:
				setVal = val.Int()
			case float64:
				setVal = val.Float64()
			case []byte:
				setVal = val.Bytes()
			case time.Time:
				setVal = val.Date()
			case bool:
				setVal = val.Bool()
			case *Delimiters:
				setVal = msg.Delimiters()
			default:
				return fmt.Errorf("hl7v2: unsupported field type %T", f.Value())
			}

			if err := f.Set(setVal); err != nil {
				return err
			}
		}
	}

	return nil
}

func Marshal(v any) ([]byte, error) {
	panic("not implemented")
}

func isSlice(v any) bool {
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array
}
