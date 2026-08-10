package hl7v2

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	msgTypeACK = "ACK"
)

// AcknowledgmentCode represents an HL7 acknowledgment code.
type AcknowledgmentCode int

const (
	// AckUnknown indicates an unknown acknowledgment code.
	AckUnknown           AcknowledgmentCode = iota //UN
	// AckApplicationAccept indicates application accept (AA).
	AckApplicationAccept                           //AA
	// AckApplicationError indicates application error (AE).
	AckApplicationError                            //AE
	// AckApplicationReject indicates application reject (AR).
	AckApplicationReject                           //AR
	// AckCommitAccept indicates commit accept (CA).
	AckCommitAccept                                //CA
	// AckCommitError indicates commit error (CE).
	AckCommitError                                 //CE
	// AckCommitReject indicates commit reject (CR).
	AckCommitReject                                //CR
)

// String returns the 2-character string representation of the AcknowledgmentCode.
func (a AcknowledgmentCode) String() string {
	switch a {
	case AckApplicationAccept:
		return "AA"
	case AckApplicationError:
		return "AE"
	case AckApplicationReject:
		return "AR"
	case AckCommitAccept:
		return "CA"
	case AckCommitError:
		return "CE"
	case AckCommitReject:
		return "CR"
	default:
		return "UN"
	}
}

// MessageErrorCode represents an HL7 error condition code used in ERR segments.
type MessageErrorCode int

const (
	// MessageErrorAccepted indicates success (0).
	MessageErrorAccepted                 MessageErrorCode = 0
	// MessageErrorSegmentSequence indicates a segment sequence error (100).
	MessageErrorSegmentSequence          MessageErrorCode = 100
	// MessageErrorRequiredFieldMissing indicates a missing required field (101).
	MessageErrorRequiredFieldMissing     MessageErrorCode = 101
	// MessageErrorDataTypeError indicates a data type error (102).
	MessageErrorDataTypeError            MessageErrorCode = 102
	// MessageErrorTableValueNotFound indicates a table value was not found (103).
	MessageErrorTableValueNotFound       MessageErrorCode = 103
	// MessageErrorUnsupportedMsgType indicates an unsupported message type (200).
	MessageErrorUnsupportedMsgType       MessageErrorCode = 200
	// MessageErrorUnsupportedEventCode indicates an unsupported event code (201).
	MessageErrorUnsupportedEventCode     MessageErrorCode = 201
	// MessageErrorUnsupportedProcessingID indicates an unsupported processing ID (202).
	MessageErrorUnsupportedProcessingID  MessageErrorCode = 202
	// MessageErrorUnsupportedVersionID indicates an unsupported version ID (203).
	MessageErrorUnsupportedVersionID     MessageErrorCode = 203
	// MessageErrorUnknownKeyID indicates an unknown key ID (204).
	MessageErrorUnknownKeyID             MessageErrorCode = 204
	// MessageErrorDuplicateKeyID indicates a duplicate key ID (205).
	MessageErrorDuplicateKeyID           MessageErrorCode = 205
	// MessageErrorAppRecordLocked indicates an application record is locked (206).
	MessageErrorAppRecordLocked          MessageErrorCode = 206
	// MessageErrorApplicationInternalError indicates an application internal error (207).
	MessageErrorApplicationInternalError MessageErrorCode = 207
)

// Acknowledgement represents an HL7 acknowledgment message structure.
type Acknowledgement struct {
	MessageHeader
	Code             string
	ControlID        string
	Message          string
	ExpectedSequence string
	DelayedAckType   string
	ErrorCondition   string
}

func buildAckHeader(msg []byte) ([]byte, []byte, *Delimiters, error) {
	if !bytes.HasPrefix(msg, []byte("MSH")) {
		return nil, nil, nil, fmt.Errorf("no MSH present")
	}

	msg = ReplaceLineEndings(msg)
	enc, err := getEncodingChars(msg)
	if err != nil {
		return nil, nil, nil, err
	}

	delims, err := ParseDelimiters(enc)
	if err != nil {
		return nil, nil, nil, err
	}

	msh := bytes.Split(msg, []byte{segmentDelimiter})[0]
	parts := delims.Split(msh, FieldDelimiter)
	if len(parts) < 10 {
		return nil, nil, nil, fmt.Errorf("invalid MSH segment: insufficient fields")
	}

	// Swap Sending/Receiving App and Facility
	parts[2], parts[4] = parts[4], parts[2] // MSH-3 and MSH-5
	parts[3], parts[5] = parts[5], parts[3] // MSH-4 and MSH-6

	// Set Message Type MSH-9 to ACK
	parts[8] = []byte(msgTypeACK)

	// Control ID MSH-10
	cid := parts[9]

	mshBytes := delims.Join(parts, FieldDelimiter)
	return mshBytes, cid, delims, nil
}

// AckRawMessage generates a raw HL7 ACK message for the given raw HL7 message.
func AckRawMessage(msg []byte) ([]byte, error) {
	mshBytes, cid, delims, err := buildAckHeader(msg)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	out.Write(mshBytes)
	out.WriteByte(segmentDelimiter)
	out.WriteString("MSA")
	out.WriteByte(delims.Field.Byte())
	out.WriteString("AA")
	out.WriteByte(delims.Field.Byte())
	out.Write(cid)

	return out.Bytes(), nil
}

// NackRawMessage generates a raw HL7 NACK message for the given raw HL7 message.
func NackRawMessage(msg []byte, code AcknowledgmentCode, errMsg string, cond ...MessageErrorCode) ([]byte, error) {
	mshBytes, cid, delims, err := buildAckHeader(msg)
	if err != nil {
		return nil, err
	}

	if code == 0 {
		code = AckApplicationAccept
	}

	condCode := ""
	if len(cond) > 0 {
		condCode = strconv.Itoa(int(cond[0]))
	}

	var out bytes.Buffer
	out.Write(mshBytes)
	out.WriteByte(segmentDelimiter)
	out.WriteString("MSA")
	out.WriteByte(delims.Field.Byte())
	out.WriteString(code.String())
	out.WriteByte(delims.Field.Byte())
	out.Write(cid)
	out.WriteByte(delims.Field.Byte())
	out.WriteString(errMsg)
	out.WriteByte(delims.Field.Byte())
	out.WriteByte(delims.Field.Byte())
	if len(cond) > 0 {
		out.WriteByte(delims.Field.Byte())
		out.WriteString(condCode)
	}

	if len(cond) > 0 || errMsg != "" {
		out.WriteByte(segmentDelimiter)
		out.WriteString("ERR")
		out.WriteByte(delims.Field.Byte())
		if len(cond) > 0 {
			out.WriteString(condCode)
			out.WriteByte(delims.Component.Byte())
			out.WriteString(errMsg)
		} else {
			out.WriteByte(delims.Component.Byte())
			out.WriteString(errMsg)
		}
	}

	return out.Bytes(), nil
}
