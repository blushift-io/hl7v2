package hl7v2

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/spf13/cast"
)

type Value struct {
	v []byte
}

func NewValue(v []byte) Value {
	return Value{
		v: v,
	}
}

func UnescapeValue(v []byte, delims *Delimiters) Value {
	if delims == nil {
		delims = DefaultDelimiters()
	}

	v = delims.Escaper().Unescape(v)

	return NewValue(v)
}

func NewStringValue(v string) Value {
	return NewValue([]byte(v))
}

func NewIntValue(v int) Value {
	return NewStringValue(strconv.Itoa(v))
}

func NewEmptyValue() Value {
	return Value{}
}

func NewDateValue(date time.Time) Value {
	if date.IsZero() {
		return NewEmptyValue()
	}

	return NewStringValue(date.Format("20060102"))
}

func NewDateTimeValue(date time.Time) Value {
	if date.IsZero() {
		return NewEmptyValue()
	}

	return NewStringValue(date.Format("20060102150405"))
}

func NewTimestampValue(t time.Time) Value {
	if t.IsZero() {
		return NewEmptyValue()
	}

	return NewStringValue(t.Format("200601021504059"))
}

func MarshalValue(v any) Value {
	if rv := reflect.ValueOf(v); rv.IsZero() || rv.IsNil() {
		return Value{}
	}

	b := cast.ToString(v)

	return Value{
		v: []byte(b),
	}
}

func (v Value) Bind(val any) error {
	rv := reflect.ValueOf(val)

	if !rv.CanSet() {
		return fmt.Errorf("hl7v2: Bind expects a settable pointer")
	}

	rv.Elem().Set(reflect.ValueOf(v.v))

	return nil
}

func (v Value) String() string {
	return string(v.v)
}

func (v Value) Bytes() []byte {
	return v.v
}

func (v Value) Int() int {
	return cast.ToInt(string(v.v))
}

func (v Value) Float64() float64 {
	return cast.ToFloat64(string(v.v))
}

func (v Value) Bool() bool {
	return cast.ToBool(string(v.v))
}

// TODO: support for timezone adjustments
func (v Value) Date() time.Time {
	if v.Empty() {
		return time.Time{}
	}

	return cast.ToTime(string(v.v))
}

func (v Value) Empty() bool {
	return len(v.v) == 0
}

func (v Value) Escape(delims *Delimiters) Value {
	if delims == nil {
		delims = DefaultDelimiters()
	}

	return NewValue(delims.Escaper().Escape(v.v))
}

func (v Value) Unescape(delims *Delimiters) Value {
	if delims == nil {
		delims = DefaultDelimiters()
	}

	return NewValue(delims.Escaper().Unescape(v.v))
}
