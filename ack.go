package hl7v2

import (
	"bytes"
	"fmt"
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

type Acknowledgement struct {
	MessageHeader
	Code             string
	ControlID        string
	Message          string
	ExpectedSequence string
	DelayedAckType   string
	ErrorCondition   string
}

func AckRawMessage(msg []byte) ([]byte, error) {
	if !bytes.HasPrefix(msg, []byte("MSH")) {
		return nil, fmt.Errorf("no MSH present")
	}

	msg = replaceLineEndings(msg)
	enc, err := getEncodingChars(msg)
	if err != nil {
		return nil, err
	}

	delims, err := ParseDelimiters(enc)
	if err != nil {
		return nil, err
	}

	msh := bytes.Split(msg, []byte{segmentDelimiter})[0]

	parts := delims.Split(msh, FieldDelimiter)
	if len(parts) < 9 {
		return nil, fmt.Errorf("invalid MSH segment")
	}

	parts[0], parts[2] = parts[2], parts[0]
	parts[1], parts[3] = parts[3], parts[1]
	parts[6] = []byte("ACK")

	cid := parts[7]
	var out bytes.Buffer

	out.Write(msh[:8])
	out.WriteByte(delims.Field.Byte())
	out.Write(delims.Join(parts, FieldDelimiter))
	out.WriteByte(segmentDelimiter)
	out.WriteString("MSA")
	out.WriteByte(delims.Field.Byte())
	out.WriteString("AA")
	out.WriteByte(delims.Field.Byte())
	out.Write(cid)

	return out.Bytes(), nil
}
