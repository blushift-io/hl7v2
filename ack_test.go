package hl7v2

import (
	"bytes"
	"testing"
)

func TestAckRawMessage(t *testing.T) {
	msg := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260101||ADT^A01|MSG123|P|2.5\rPID|1||12345\r")

	ackBytes, err := AckRawMessage(msg)
	if err != nil {
		t.Fatalf("failed to generate ACK: %v", err)
	}

	if !bytes.Contains(ackBytes, []byte("MSA|AA|MSG123")) {
		t.Errorf("expected MSA|AA|MSG123 in ACK, got: %s", string(ackBytes))
	}
}

func TestNackRawMessage(t *testing.T) {
	msg := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260101||ADT^A01|MSG123|P|2.5\rPID|1||12345\r")

	nackBytes, err := NackRawMessage(msg, AckApplicationError, "Application Error", MessageErrorApplicationInternalError)
	if err != nil {
		t.Fatalf("failed to generate NACK: %v", err)
	}

	nackStr := string(nackBytes)
	if !bytes.Contains(nackBytes, []byte("MSA|AE|MSG123|Application Error|||207")) {
		t.Errorf("expected MSA|AE|MSG123|Application Error|||207 in NACK, got: %s", nackStr)
	}

	if !bytes.Contains(nackBytes, []byte("ERR|207^Application Error")) {
		t.Errorf("expected ERR|207^Application Error in NACK, got: %s", nackStr)
	}
}
