package gen

import (
	"fmt"
	"reflect"
	"time"

	"math/rand"
)

func Generate(v any) error {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("generate input must be a pointer to a struct")
	}

	if reflect.ValueOf(v).IsNil() {
		return fmt.Errorf("generate input is nil")
	}

	rv := reflect.ValueOf(v)
	val, err := generateValues(v)
	if err != nil {
		return err
	}

	if rv.Elem().CanSet() && rv.Elem().CanConvert(rt.Elem()) {
		rv.Elem().Set(val.Elem().Convert(rt.Elem()))
	}

	return nil
}

func generateValues(v any) (reflect.Value, error) {
	t := reflect.TypeOf(v)
	if t == nil {
		return reflect.Value{}, fmt.Errorf("generate input is nil")
	}

	k := t.Kind()
	switch k {
	case reflect.Ptr:
		nv := reflect.New(t.Elem())

		var res reflect.Value
		var err error

		if v != reflect.Zero(reflect.TypeOf(v)).Interface() {
			res, err = generateValues(reflect.ValueOf(v).Elem().Interface())
		} else {
			res, err = generateValues(nv.Elem().Interface())
		}

		if err != nil {
			return reflect.Value{}, err
		}

		if reflect.ValueOf(nv).IsZero() {
			return res, nil
		}

		if nv.Elem().CanSet() && nv.Elem().CanConvert(t.Elem()) {
			nv.Elem().Set(res.Convert(nv.Elem().Type()))
		}

		return nv, nil
	case reflect.Struct:
		switch t.String() {
		case "time.Time":
			tt := time.Now().Add(time.Duration(rand.Int63()))
			return reflect.ValueOf(tt), nil
		default:
			nv := reflect.New(t).Elem()

			for i := 0; i < nv.NumField(); i++ {
				f := nv.Field(i)

				if !f.CanSet() {
					continue
				}

				tag := t.Field(i).Tag.Get("gen")
				if len(tag) == 0 {
					continue
				}

				if err := setTaggedField(f.Addr(), tag); err != nil {
					return reflect.Value{}, err
				}
			}
		}
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type %s", k)
	}

	return reflect.Value{}, nil
}

func setTaggedField(v reflect.Value, tag string) error {
	if f, ok := genFuncs[tag]; ok {
		return f(v)
	}

	return fmt.Errorf("no generator function found for tag %s", tag)
}
