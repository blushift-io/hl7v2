package hl7v2

import (
	"fmt"
	"reflect"
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

func NewValueString(v string) Value {
	return NewValue([]byte(v))
}

func NewEmptyValue() Value {
	return Value{}
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
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("hl7v2: Bind expects a pointer")
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

func (v Value) Date() time.Time {
	return cast.ToTime(string(v.v))
}

func (v Value) Empty() bool {
	return v.v == nil
}
