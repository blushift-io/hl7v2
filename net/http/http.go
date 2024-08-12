package http

import (
	"encoding/base64"

	"github.com/blushift-io/hl7v2"
)

type Message struct {
	Data string `json:"data"`
}

func NewMessage(data []byte) *Message {
	enc := base64.StdEncoding.EncodeToString(data)

	return &Message{Data: enc}
}

func (m Message) RawMessage() (*hl7v2.RawMessage, error) {
	if len(m.Data) == 0 {
		return nil, nil
	}

	dec, err := base64.StdEncoding.DecodeString(m.Data)
	if err != nil {
		return nil, err
	}

	return hl7v2.ParseRaw(dec)
}
