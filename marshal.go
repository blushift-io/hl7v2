package hl7v2

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func marshalHL7(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	if m, ok := v.(Marshaler); ok {
		return m.MarshalHL7()
	}

	delims := DefaultDelimiters()
	enc := newEncoder(delims)
	return enc.encodeTop(v)
}

type encoder struct {
	delims *Delimiters
}

func newEncoder(delims *Delimiters) *encoder {
	if delims == nil {
		delims = DefaultDelimiters()
	}
	return &encoder{delims: delims}
}

func (enc *encoder) encodeTop(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)
	rv = unwrapValue(rv)
	if !rv.IsValid() {
		return nil, nil
	}

	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return rv.Bytes(), nil
		}
		var segLines [][]byte
		for i := 0; i < rv.Len(); i++ {
			b, err := enc.encodeTop(rv.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			if len(b) > 0 {
				segLines = append(segLines, b)
			}
		}
		return bytes.Join(segLines, nil), nil
	}

	if rv.Kind() == reflect.Struct {
		if isMessageOrGroupStruct(rv.Type()) {
			return enc.encodeMessageOrGroup(rv)
		}
		return enc.encodeSegment(rv, "")
	}

	return enc.encodePrimitive(rv), nil
}

func unwrapValue(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return reflect.Value{}
		}
		rv = rv.Elem()
	}
	return rv
}

func isMessageOrGroupStruct(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		tag := sf.Tag.Get("hl7")
		if tag == "" {
			continue
		}
		tagKey := strings.Split(tag, ",")[0]
		if tagKey != "" && !isNumeric(tagKey) {
			return true
		}
	}
	return false
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (enc *encoder) encodeMessageOrGroup(rv reflect.Value) ([]byte, error) {
	t := rv.Type()
	var segLines [][]byte

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		tag := sf.Tag.Get("hl7")
		if tag == "" {
			continue
		}
		segName := strings.Split(tag, ",")[0]
		fv := unwrapValue(rv.Field(i))
		if !fv.IsValid() {
			continue
		}

		if fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array {
			for j := 0; j < fv.Len(); j++ {
				item := unwrapValue(fv.Index(j))
				if !item.IsValid() {
					continue
				}
				b, err := enc.encodeItemInMessage(item, segName)
				if err != nil {
					return nil, err
				}
				if len(b) > 0 {
					segLines = append(segLines, b)
				}
			}
		} else {
			b, err := enc.encodeItemInMessage(fv, segName)
			if err != nil {
				return nil, err
			}
			if len(b) > 0 {
				segLines = append(segLines, b)
			}
		}
	}

	return bytes.Join(segLines, nil), nil
}

func (enc *encoder) encodeItemInMessage(fv reflect.Value, segName string) ([]byte, error) {
	if fv.Kind() == reflect.Struct {
		if isMessageOrGroupStruct(fv.Type()) {
			return enc.encodeMessageOrGroup(fv)
		}
		return enc.encodeSegment(fv, segName)
	}
	return nil, nil
}

func (enc *encoder) encodeSegment(rv reflect.Value, overrideSegID string) ([]byte, error) {
	t := rv.Type()
	segID := overrideSegID
	if segID == "" {
		segID = t.Name()
	}

	fieldMap := make(map[int]reflect.Value)
	maxPos := 0

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		tag := sf.Tag.Get("hl7")
		if tag == "" {
			continue
		}
		posStr := strings.Split(tag, ",")[0]
		pos, err := strconv.Atoi(posStr)
		if err != nil {
			continue
		}
		fv := unwrapValue(rv.Field(i))
		if fv.IsValid() {
			fieldMap[pos] = fv
			if pos > maxPos {
				maxPos = pos
			}
		}
	}

	if segID == "MSH" {
		if maxPos < 2 {
			maxPos = 2
		}
	}

	var fieldBytes [][]byte
	for pos := 1; pos <= maxPos; pos++ {
		fv, ok := fieldMap[pos]
		if segID == "MSH" && pos == 1 {
			b := enc.delims.Field.Bytes()
			if ok && fv.IsValid() && fv.Kind() == reflect.String && fv.String() != "" {
				b = []byte(fv.String())
			}
			fieldBytes = append(fieldBytes, b)
			continue
		}
		if segID == "MSH" && pos == 2 {
			b := enc.delims.EncodingCharsValue().Bytes()
			if ok && fv.IsValid() && fv.Kind() == reflect.String && fv.String() != "" {
				b = []byte(fv.String())
			}
			fieldBytes = append(fieldBytes, b)
			continue
		}

		if !ok || !fv.IsValid() {
			fieldBytes = append(fieldBytes, nil)
		} else {
			fb := enc.encodeFieldVal(fv)
			fieldBytes = append(fieldBytes, fb)
		}
	}

	var buf bytes.Buffer
	if segID == "MSH" {
		buf.WriteString("MSH")
		if len(fieldBytes) >= 1 {
			buf.Write(fieldBytes[0])
		} else {
			buf.Write(enc.delims.Field.Bytes())
		}
		if len(fieldBytes) >= 2 {
			buf.Write(fieldBytes[1])
		} else {
			buf.Write(enc.delims.EncodingCharsValue().Bytes())
		}
		buf.Write(enc.delims.Field.Bytes())

		rest := fieldBytes
		if len(rest) >= 2 {
			rest = rest[2:]
		} else {
			rest = nil
		}
		rest = trimTrailingEmpty(rest)
		for i, f := range rest {
			if i > 0 {
				buf.Write(enc.delims.Field.Bytes())
			}
			buf.Write(f)
		}
	} else {
		buf.WriteString(segID)
		trimmed := trimTrailingEmpty(fieldBytes)
		for _, f := range trimmed {
			buf.Write(enc.delims.Field.Bytes())
			buf.Write(f)
		}
	}

	buf.Write(enc.delims.Segment.Bytes())
	return buf.Bytes(), nil
}

func (enc *encoder) encodeFieldVal(fv reflect.Value) []byte {
	if fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array {
		if fv.Type().Elem().Kind() == reflect.Uint8 {
			return fv.Bytes()
		}
		var reps [][]byte
		for i := 0; i < fv.Len(); i++ {
			item := unwrapValue(fv.Index(i))
			cb := enc.encodeComponentVal(item)
			reps = append(reps, cb)
		}
		reps = trimTrailingEmpty(reps)
		return bytes.Join(reps, enc.delims.Repetition.Bytes())
	}

	return enc.encodeComponentVal(fv)
}

func (enc *encoder) encodeComponentVal(cv reflect.Value) []byte {
	if !cv.IsValid() {
		return nil
	}
	if cv.Kind() == reflect.Struct {
		return enc.encodeComponentStruct(cv)
	}
	return enc.encodePrimitive(cv)
}

func (enc *encoder) encodeComponentStruct(rv reflect.Value) []byte {
	t := rv.Type()
	compMap := make(map[int]reflect.Value)
	maxPos := 0

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		tag := sf.Tag.Get("hl7")
		if tag == "" {
			continue
		}
		posStr := strings.Split(tag, ",")[0]
		pos, err := strconv.Atoi(posStr)
		if err != nil {
			continue
		}
		fv := unwrapValue(rv.Field(i))
		if fv.IsValid() {
			compMap[pos] = fv
			if pos > maxPos {
				maxPos = pos
			}
		}
	}

	var comps [][]byte
	for pos := 1; pos <= maxPos; pos++ {
		fv, ok := compMap[pos]
		if !ok || !fv.IsValid() {
			comps = append(comps, nil)
		} else {
			comps = append(comps, enc.encodeSubcomponentVal(fv))
		}
	}

	comps = trimTrailingEmpty(comps)
	return bytes.Join(comps, enc.delims.Component.Bytes())
}

func (enc *encoder) encodeSubcomponentVal(sv reflect.Value) []byte {
	if !sv.IsValid() {
		return nil
	}
	if sv.Kind() == reflect.Struct {
		return enc.encodeSubcomponentStruct(sv)
	}
	return enc.encodePrimitive(sv)
}

func (enc *encoder) encodeSubcomponentStruct(rv reflect.Value) []byte {
	t := rv.Type()
	subMap := make(map[int]reflect.Value)
	maxPos := 0

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		tag := sf.Tag.Get("hl7")
		if tag == "" {
			continue
		}
		posStr := strings.Split(tag, ",")[0]
		pos, err := strconv.Atoi(posStr)
		if err != nil {
			continue
		}
		fv := unwrapValue(rv.Field(i))
		if fv.IsValid() {
			subMap[pos] = fv
			if pos > maxPos {
				maxPos = pos
			}
		}
	}

	var subs [][]byte
	for pos := 1; pos <= maxPos; pos++ {
		fv, ok := subMap[pos]
		if !ok || !fv.IsValid() {
			subs = append(subs, nil)
		} else {
			subs = append(subs, enc.encodePrimitive(fv))
		}
	}

	subs = trimTrailingEmpty(subs)
	return bytes.Join(subs, enc.delims.Subcomponent.Bytes())
}

func (enc *encoder) encodePrimitive(v reflect.Value) []byte {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		return []byte(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		return []byte(strconv.FormatFloat(v.Float(), 'f', -1, 64))
	case reflect.Bool:
		if v.Bool() {
			return []byte("Y")
		}
		return []byte("N")
	default:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			t := v.Interface().(time.Time)
			if t.IsZero() {
				return nil
			}
			if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
				return []byte(t.Format("20060102"))
			}
			return []byte(t.Format("20060102150405"))
		}
		return []byte(fmt.Sprint(v.Interface()))
	}
}

func trimTrailingEmpty(items [][]byte) [][]byte {
	last := len(items) - 1
	for last >= 0 && len(items[last]) == 0 {
		last--
	}
	if last < 0 {
		return nil
	}
	return items[:last+1]
}
