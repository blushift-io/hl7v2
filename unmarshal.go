package hl7v2

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Unmarshal unmarshals HL7 wire format bytes into the provided struct pointer or Unmarshaler.
func Unmarshal(b []byte, v any) error {
	if v == nil {
		return fmt.Errorf("hl7v2: Unmarshal target cannot be nil")
	}

	if u, ok := v.(Unmarshaler); ok {
		return u.UnmarshalHL7(b)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("hl7v2: Unmarshal expects a non-nil pointer to a struct")
	}

	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("hl7v2: Unmarshal expects a pointer to a struct")
	}

	raw, err := ParseRaw(b)
	if err != nil {
		return fmt.Errorf("hl7v2: failed to parse raw message: %w", err)
	}

	return unmarshalStruct(raw, elem, "")
}

func unmarshalStruct(raw *RawMessage, rv reflect.Value, prefix string) error {
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil
	}

	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}

		tag := sf.Tag.Get("hl7")
		if len(tag) == 0 {
			continue
		}

		tagOpts := strings.Split(tag, ",")
		var queryLocs []string
		for _, opt := range tagOpts {
			q := strings.TrimSpace(opt)
			if q == "required" || q == "repeatable" {
				continue
			}
			queryLocs = append(queryLocs, q)
		}

		fieldVal := rv.Field(i)

		for _, tagKey := range queryLocs {
			var queryLoc string
			if prefix == "" {
				queryLoc = tagKey
			} else {
				if isNumeric(tagKey) {
					queryLoc = fmt.Sprintf("%s.%s", prefix, tagKey)
				} else {
					queryLoc = tagKey
				}
			}

			targetKind := sf.Type.Kind()
			if targetKind == reflect.Pointer {
				targetKind = sf.Type.Elem().Kind()
			}

			if targetKind == reflect.Struct && !isSpecialFieldType(sf.Type) {
				if err := unmarshalStruct(raw, fieldVal, queryLoc); err != nil {
					return err
				}
				continue
			}

			val, err := raw.QueryValue(queryLoc)
			if err != nil {
				if errors.Is(err, ErrValueNotFound) || errors.Is(err, ErrElementNotFound) {
					continue
				}
				return err
			}

			if val == nil || val.Empty() {
				continue
			}

			if err := setFieldValue(fieldVal, val, raw); err != nil {
				return err
			}
		}
	}

	return nil
}

func isSpecialFieldType(t reflect.Type) bool {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t == reflect.TypeOf(time.Time{}) || t == reflect.TypeOf(Delimiters{}) || t == reflect.TypeOf(&Delimiters{})
}

func setFieldValue(f reflect.Value, val *Value, raw *RawMessage) error {
	if !f.IsValid() || !f.CanSet() {
		return nil
	}

	switch f.Kind() {
	case reflect.String:
		f.SetString(val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.SetInt(int64(val.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f.SetUint(uint64(val.Int()))
	case reflect.Float32, reflect.Float64:
		f.SetFloat(val.Float64())
	case reflect.Bool:
		f.SetBool(val.Bool())
	case reflect.Slice:
		if f.Type().Elem().Kind() == reflect.Uint8 {
			f.SetBytes(val.Bytes())
		}
	default:
		if f.Type() == reflect.TypeOf(time.Time{}) {
			f.Set(reflect.ValueOf(val.Date()))
		} else if f.Type() == reflect.TypeOf(&Delimiters{}) {
			f.Set(reflect.ValueOf(raw.Delimiters()))
		} else if f.Type() == reflect.TypeOf(Delimiters{}) {
			f.Set(reflect.ValueOf(*raw.Delimiters()))
		} else {
			return fmt.Errorf("hl7v2: unsupported field type %s", f.Type())
		}
	}
	return nil
}
