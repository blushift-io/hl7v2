package hl7v2

import "errors"

var (
	ErrMsgLength = errors.New("invalid message: message length must be at least 8 bytes")
)
