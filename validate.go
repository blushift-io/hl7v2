package hl7v2

import (
	"fmt"
	"strings"

	"github.com/blushift-io/hl7v2/schema"
)

// ValidationResult holds the result of a validation operation and any validation errors.
type ValidationResult struct {
	Valid  bool
	Errors []error
}

// Validator is the interface implemented by types that can validate an HL7 Element.
type Validator interface {
	Validate(Element) ValidationResult
}

// ValidationFunc is an adapter allowing a function to be used as a Validator.
type ValidationFunc func(Element) ValidationResult

// Validate implements the Validator interface for ValidationFunc.
func (f ValidationFunc) Validate(e Element) ValidationResult {
	return f(e)
}

// ValidateMessageSchema validates an HL7 message against its schema specification.
func ValidateMessageSchema(el Element) ValidationResult {
	if el.Type() != ElementMessage {
		return ValidationResult{false, []error{fmt.Errorf("invalid element type: %s", el.Type())}}
	}

	msg, ok := el.(*Message)
	if !ok {
		return ValidationResult{false, []error{fmt.Errorf("invalid message type: %T", el)}}
	}

	msgTyp := msg.Header().MessageType().String()

	sch, err := msg.Schema()
	if err != nil {
		return ValidationResult{false, []error{err}}
	}

	for _, seg := range sch.Segments {
		if !seg.Group && seg.Required() && !msg.HasSegment(seg.ID) {
			return ValidationResult{false, []error{fmt.Errorf("missing required segment: %s", seg.ID)}}
		}
	}

	for _, seg := range msg.Segments() {
		if strings.HasPrefix(seg.Name(), "Z") {
			continue
		}

		msgSeg := sch.Segment(seg.Name())
		if msgSeg == nil {
			return ValidationResult{false, []error{fmt.Errorf("unexpected segment %s in message type %s", seg.Name(), msgTyp)}}
		}

		schSeg := msgSeg.Segment()
		if schSeg == nil {
			return ValidationResult{false, []error{fmt.Errorf("segment %s not found in v%s schema", seg.Name(), sch.Version())}}
		}

		res := validateSegmentSchema(seg, schSeg)
		if !res.Valid {
			return res
		}
	}

	return ValidationResult{true, nil}
}

// ValidateSegmentSchema validates an HL7 segment against its schema specification.
func ValidateSegmentSchema(el Element) ValidationResult {
	seg, ok := el.(*Segment)
	if !ok {
		return ValidationResult{false, []error{fmt.Errorf("invalid segment type: %T", el)}}
	}

	h := seg.Header()
	if h == nil {
		return ValidationResult{false, []error{fmt.Errorf("missing segment header")}}
	}

	sch := schema.Open(h.VersionID)
	if sch == nil {
		return ValidationResult{false, []error{fmt.Errorf("schema version %s not found", h.VersionID)}}
	}

	schSeg := sch.Segment(seg.Name())
	if schSeg == nil {
		return ValidationResult{false, []error{fmt.Errorf("segment %s not found in v%s schema", seg.Name(), sch.Version())}}
	}

	return validateSegmentSchema(seg, schSeg)
}

func validateSegmentSchema(seg *Segment, sch *schema.Segment) ValidationResult {
	for i, sf := range sch.Fields {
		f, err := seg.Field(i + 1)
		if err != nil {
			if sf.Required() {
				return ValidationResult{false, []error{fmt.Errorf("missing required field %d: %s", i, sf.ID)}}
			}

			continue
		}

		if sf.MaxRepetitions > 0 && f.Length() > sf.MaxRepetitions {
			return ValidationResult{false, []error{fmt.Errorf("field %s (%d) exceeds max repetitions: %d", sf.ID, f.Length(), sf.MaxRepetitions)}}
		}

		if sf.MaxRepetitions > 0 && f.Length() < sf.MinRepetitions {
			return ValidationResult{false, []error{fmt.Errorf("field %s below min repetitions: %d", sf.ID, sf.MinRepetitions)}}
		}

		if sf.Length > 0 {
			if v := f.Value().String(); len(v) > sf.Length {
				return ValidationResult{false, []error{fmt.Errorf("field %s (%d) exceeds max length: %d", sf.ID, len(v), sf.Length)}}
			}
		}

		if sf.Table() != nil {
			if err := validateTableValue(f.Value(), sf.Table()); err != nil {
				return ValidationResult{false, []error{fmt.Errorf("field %s has invalid table value: %s", sf.ID, err)}}
			}
		}
	}

	return ValidationResult{true, nil}
}

func validateTableValue(v Value, tbl *schema.Table) error {
	if tbl == nil {
		return fmt.Errorf("table is nil")
	}

	switch tbl.Type {
	case schema.TableTypeHL7:
		for _, tv := range tbl.Entries {
			if v.String() == tv.Value {
				return nil
			}
		}
	case schema.TableTypeLocal, schema.TableTypeUnknown, schema.TableTypeUser:
		return nil
	}

	return fmt.Errorf("invalid table value: %s", v.String())
}
