package hl7v2

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/spf13/cast"
)

// Value represents an HL7 element value as raw bytes with helper conversion methods.
type Value struct {
	v []byte
}

// NewValue creates a new Value wrapping the provided byte slice.
func NewValue(v []byte) Value {
	return Value{
		v: v,
	}
}

// UnescapeValue creates a new Value from unescaping the input byte slice with the given delimiters.
func UnescapeValue(v []byte, delims *Delimiters) Value {
	if delims == nil {
		delims = DefaultDelimiters()
	}

	v = delims.Escaper().Unescape(v)

	return NewValue(v)
}

// NewStringValue creates a new Value from a string.
func NewStringValue(v string) Value {
	return NewValue([]byte(v))
}

// NewIntValue creates a new Value from an integer.
func NewIntValue(v int) Value {
	return NewStringValue(strconv.Itoa(v))
}

// NewEmptyValue returns an empty Value.
func NewEmptyValue() Value {
	return Value{}
}

// NewDateValue creates a new Value formatted as an HL7 date (YYYYMMDD) from a time.Time.
func NewDateValue(date time.Time) Value {
	if date.IsZero() {
		return NewEmptyValue()
	}

	return NewStringValue(date.Format("20060102"))
}

// NewDateTimeValue creates a new Value formatted as an HL7 datetime (YYYYMMDDHHMMSS) from a time.Time.
func NewDateTimeValue(date time.Time) Value {
	if date.IsZero() {
		return NewEmptyValue()
	}

	return NewStringValue(date.Format("20060102150405"))
}

// NewTimestampValue creates a new Value formatted as an HL7 timestamp from a time.Time.
func NewTimestampValue(t time.Time) Value {
	if t.IsZero() {
		return NewEmptyValue()
	}

	return NewStringValue(t.Format("200601021504059"))
}

// MarshalValue creates a Value by casting the input to a string.
func MarshalValue(v any) Value {
	if rv := reflect.ValueOf(v); rv.IsZero() || rv.IsNil() {
		return Value{}
	}

	b := cast.ToString(v)

	return Value{
		v: []byte(b),
	}
}

// Bind binds the byte content of Value to a target variable pointer.
func (v Value) Bind(val any) error {
	rv := reflect.ValueOf(val)

	if !rv.CanSet() {
		return fmt.Errorf("hl7v2: Bind expects a settable pointer")
	}

	rv.Elem().Set(reflect.ValueOf(v.v))

	return nil
}

// String returns the Value as a string.
func (v Value) String() string {
	return string(v.v)
}

// Bytes returns the underlying raw byte slice of the Value.
func (v Value) Bytes() []byte {
	return v.v
}

// Int converts the Value to an integer.
func (v Value) Int() int {
	return cast.ToInt(string(v.v))
}

// Float64 converts the Value to a float64.
func (v Value) Float64() float64 {
	return cast.ToFloat64(string(v.v))
}

// Bool converts the Value to a boolean.
func (v Value) Bool() bool {
	return cast.ToBool(string(v.v))
}

// Date converts the Value to a time.Time object.
// TODO: support for timezone adjustments
func (v Value) Date() time.Time {
	if v.Empty() {
		return time.Time{}
	}

	return cast.ToTime(string(v.v))
}

// Empty reports whether the Value contains no bytes.
func (v Value) Empty() bool {
	return len(v.v) == 0
}

// Escape returns a new Value with HL7 delimiters escaped.
func (v Value) Escape(delims *Delimiters) Value {
	if delims == nil {
		delims = DefaultDelimiters()
	}

	return NewValue(delims.Escaper().Escape(v.v))
}

// Unescape returns a new Value with HL7 escape sequences unescaped.
func (v Value) Unescape(delims *Delimiters) Value {
	if delims == nil {
		delims = DefaultDelimiters()
	}

	return NewValue(delims.Escaper().Unescape(v.v))
}
