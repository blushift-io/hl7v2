package hl7v2

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	msgTypeACK = "ACK"
)

type AcknowledgmentCode int

const (
	AckUnknown           AcknowledgmentCode = iota //UN
	AckApplicationAccept                           //AA
	AckApplicationError                            //AE
	AckApplicationReject                           //AR
	AckCommitAccept                                //CA
	AckCommitError                                 //CE
	AckCommitReject                                //CR
)

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

type MessageErrorCode int

const (
	MessageErrorAccepted                 MessageErrorCode = 0
	MessageErrorSegmentSequence          MessageErrorCode = 100
	MessageErrorRequiredFieldMissing     MessageErrorCode = 101
	MessageErrorDataTypeError            MessageErrorCode = 102
	MessageErrorTableValueNotFound       MessageErrorCode = 103
	MessageErrorUnsupportedMsgType       MessageErrorCode = 200
	MessageErrorUnsupportedEventCode     MessageErrorCode = 201
	MessageErrorUnsupportedProcessingID  MessageErrorCode = 202
	MessageErrorUnsupportedVersionID     MessageErrorCode = 203
	MessageErrorUnknownKeyID             MessageErrorCode = 204
	MessageErrorDuplicateKeyID           MessageErrorCode = 205
	MessageErrorAppRecordLocked          MessageErrorCode = 206
	MessageErrorApplicationInternalError MessageErrorCode = 207
)

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
