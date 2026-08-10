package hl7v2

import "errors"

var (
	// ErrMsgLength is returned when an HL7 message is shorter than the minimum length of 8 bytes.
	ErrMsgLength       = errors.New("invalid message: message length must be at least 8 bytes")
	// ErrElementNotFound is returned when a requested element is not found.
	ErrElementNotFound = errors.New("element not found")
	// ErrValueNotFound is returned when a requested value is not found.
	ErrValueNotFound   = errors.New("value not found")
)
