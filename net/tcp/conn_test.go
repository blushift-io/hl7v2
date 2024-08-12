package tcp

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/blushift-io/hl7v2"
)

func TestConnDial(t *testing.T) {
	conn, err := Dial(context.Background(), "localhost:5050")
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	defer conn.Close()

	b, err := os.ReadFile("../../test/phi/852855.hl7")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	msg, err := hl7v2.ParseRaw(b, hl7v2.FixLineEndings())
	if err != nil {
		t.Fatalf("failed to parse message: %v", err)
	}

	if err := conn.WriteMessage(msg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	ack, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	fmt.Println(string(ack.Value().Bytes()))
}
