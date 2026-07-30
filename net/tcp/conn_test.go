package tcp

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/blushift-io/hl7v2"
)

func TestConnDial(t *testing.T) {
	conn, err := Dial(context.Background(), "127.0.0.1:5050", WithConnOption(DefaultConnOptions()))
	if err != nil {
		t.Skipf("skipping integration dial test: %v", err)
	}

	defer conn.Close()

	b, err := os.ReadFile("../../test/phi/852855.hl7")
	if err != nil {
		t.Skipf("failed to read test file: %v", err)
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

func WithConnOption(opts *ConnOptions) ConnOption {
	return func(o *ConnOptions) {
		*o = *opts
		o.DialRetries = 0
	}
}
