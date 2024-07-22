package hl7v2

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNewMessage(t *testing.T) {

	msg, err := NewMessageFromFile("./fixtures/gcp_example.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	el, err := msg.Query("OBX[1].3.1")
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(el.Value())
}
