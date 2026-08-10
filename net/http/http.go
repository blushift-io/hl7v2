// Package http provides HTTP handlers and utilities for transmitting and receiving
// HL7 v2 messages as JSON or ER7 payloads.
package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/blushift-io/hl7v2"
)

// Message represents a JSON-encapsulated HL7 v2 message payload.
type Message struct {
	Data string `json:"data"`
}

// NewMessage creates a Message containing base64-encoded binary data.
func NewMessage(data []byte) *Message {
	enc := base64.StdEncoding.EncodeToString(data)

	return &Message{Data: enc}
}

// RawMessage decodes the base64 payload and parses it into a RawMessage.
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

// Handle wraps a message handling function into an http.HandlerFunc supporting ER7 and JSON content types.
func Handle(fn func(ctx context.Context, msg *hl7v2.RawMessage) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if strings.Contains(ct, ";") {
			ct = strings.TrimSpace(strings.Split(ct, ";")[0])
		}

		switch ct {
		case "x-application/hl7-v2+er7":
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			defer r.Body.Close()

			msg, err := hl7v2.ParseRaw(b)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := fn(r.Context(), msg); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		case "application/json", "x-application/hl7-v2+json":
			var m Message
			if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			msg, err := m.RawMessage()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := fn(r.Context(), msg); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		}
	})
}

// WriteResponse writes an HL7 v2 RawMessage to an http.ResponseWriter using ER7 or JSON formatting.
func WriteResponse(w http.ResponseWriter, contentType string, msg *hl7v2.RawMessage) error {
	switch contentType {
	case "x-application/hl7-v2+er7":
		w.Header().Set("Content-Type", "x-application/hl7-v2+er7")
		_, err := w.Write(msg.Value().Bytes())

		return err
	default:
		m := NewMessage(msg.Value().Bytes())
		w.Header().Set("Content-Type", "application/json")

		return json.NewEncoder(w).Encode(m)
	}
}
