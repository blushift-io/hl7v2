package hl7v2

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRawMessage(t *testing.T) {
	b, err := os.ReadFile("./test/fixtures/ORU_R01_1.hl7")
	if err != nil {
		t.Fatal(err)
	}

	m, err := ParseRaw(b, FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	j, err := m.JSON(true)
	if err != nil {
		t.Fatal(err)
	}

	if len(j) == 0 {
		t.Fatal("expected non-empty JSON output")
	}
}

func TestQueryRawMessage(t *testing.T) {
	b, err := os.ReadFile("./test/fixtures/ADT_A01_2.hl7")
	if err != nil {
		t.Fatal(err)
	}

	m, err := ParseRaw(b, FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	v, err := m.QueryValue("PV1.3[0].1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "2000", v.String())
}

func TestRawMessage_MarshalJSON(t *testing.T) {
	er7 := []byte("MSH|^~\\&|SENDING_APP|SENDING_FAC|REC_APP|REC_FAC|20260730120000||ADT^A01^ADT_A01|MSG00001|P|2.5\rPID|1||10001^^^HOSPITAL||EVERYMAN^ADAM\r")
	m, err := ParseRaw(er7)
	if err != nil {
		t.Fatalf("failed to parse raw message: %v", err)
	}

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal RawMessage to JSON: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("failed to unmarshal generated JSON: %v", err)
	}

	if parsed["delimiters"] != "|^~\\&" {
		t.Errorf("expected delimiters '|^~\\&', got %v", parsed["delimiters"])
	}

	segs, ok := parsed["segments"].([]any)
	if !ok || len(segs) != 2 {
		t.Fatalf("expected 2 segments, got %v", parsed["segments"])
	}

	msh := segs[0].(map[string]any)
	if msh["id"] != "MSH" {
		t.Errorf("expected MSH segment id, got %v", msh["id"])
	}

	pid := segs[1].(map[string]any)
	if pid["id"] != "PID" {
		t.Errorf("expected PID segment id, got %v", pid["id"])
	}

	fields := pid["fields"].(map[string]any)
	if fields["1"] != "1" {
		t.Errorf("expected PID.1 to be '1', got %v", fields["1"])
	}
	if fields["3.1"] != "10001" {
		t.Errorf("expected PID.3.1 to be '10001', got %v", fields["3.1"])
	}
	if fields["3.4"] != "HOSPITAL" {
		t.Errorf("expected PID.3.4 to be 'HOSPITAL', got %v", fields["3.4"])
	}
	if fields["5.1"] != "EVERYMAN" {
		t.Errorf("expected PID.5.1 to be 'EVERYMAN', got %v", fields["5.1"])
	}
	if fields["5.2"] != "ADAM" {
		t.Errorf("expected PID.5.2 to be 'ADAM', got %v", fields["5.2"])
	}
}

func TestRawMessage_UnmarshalJSON(t *testing.T) {
	jsonInput := `{
		"delimiters": "|^~\\&",
		"segments": [
			{
				"id": "MSH",
				"fields": {
					"1": "|",
					"2": "^~\\&",
					"3": "MY_APP",
					"4": "MY_FAC",
					"9.1": "ADT",
					"9.2": "A01"
				}
			},
			{
				"id": "PID",
				"fields": {
					"1": "1",
					"3[0].1": "12345",
					"3[0].4": "MRN",
					"3[1].1": "67890",
					"3[1].4": "SSN",
					"5.1": "DOE",
					"5.2": "JOHN"
				}
			}
		]
	}`

	var m RawMessage
	if err := json.Unmarshal([]byte(jsonInput), &m); err != nil {
		t.Fatalf("failed to unmarshal JSON into RawMessage: %v", err)
	}

	if len(m.Segments()) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(m.Segments()))
	}

	pidSegs := m.Segments("PID")
	if len(pidSegs) != 1 {
		t.Fatalf("expected 1 PID segment, got %d", len(pidSegs))
	}

	val, err := m.QueryValue("PID.1")
	if err != nil || val.String() != "1" {
		t.Errorf("expected PID.1 to be '1', got %v (err: %v)", val, err)
	}

	val, err = m.QueryValue("PID.3[0].1")
	if err != nil || val.String() != "12345" {
		t.Errorf("expected PID.3[0].1 to be '12345', got %v (err: %v)", val, err)
	}

	val, err = m.QueryValue("PID.3[1].4")
	if err != nil || val.String() != "SSN" {
		t.Errorf("expected PID.3[1].4 to be 'SSN', got %v (err: %v)", val, err)
	}
}

func TestRawMessage_JSONRoundTrip(t *testing.T) {
	b, err := os.ReadFile("./test/fixtures/ADT_A01_1.hl7")
	if err != nil {
		t.Fatal(err)
	}

	m1, err := ParseRaw(b, FixLineEndings())
	if err != nil {
		t.Fatalf("failed to parse initial raw message: %v", err)
	}

	jsonBytes, err := m1.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to serialize to JSON: %v", err)
	}

	m2, err := ParseJSON(jsonBytes)
	if err != nil {
		t.Fatalf("failed to deserialize from JSON: %v", err)
	}

	if len(m1.Segments()) != len(m2.Segments()) {
		t.Fatalf("segment count mismatch: m1 has %d, m2 has %d", len(m1.Segments()), len(m2.Segments()))
	}

	for i, s1 := range m1.Segments() {
		s2 := m2.Segments()[i]
		if s1.ID() != s2.ID() {
			t.Errorf("segment %d ID mismatch: %s vs %s", i, s1.ID(), s2.ID())
		}
	}

	// Verify query values match between original and deserialized
	v1, err1 := m1.QueryValue("PID.5.1")
	v2, err2 := m2.QueryValue("PID.5.1")
	if err1 != nil || err2 != nil || v1.String() != v2.String() {
		t.Errorf("query PID.5.1 mismatch: %v (err: %v) vs %v (err: %v)", v1, err1, v2, err2)
	}

	// Re-serialize m2 and verify JSON matches m1's JSON
	jsonBytes2, err := m2.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to re-serialize m2 to JSON: %v", err)
	}

	if !bytes.Equal(jsonBytes, jsonBytes2) {
		t.Errorf("round trip JSON mismatch:\nm1: %s\nm2: %s", string(jsonBytes), string(jsonBytes2))
	}
}

func TestRawMessage_SubcomponentsAndCustomDelimiters(t *testing.T) {
	er7 := []byte("MSH#^~\\&#MYAPP#MYFAC\rPR1#1#SUB1&SUB2^COMP2\r")
	delims, err := ParseDelimiters([]byte("#^~\\&"))
	if err != nil {
		t.Fatal(err)
	}

	m1, err := ParseRaw(er7, PreParse(func(b []byte) ([]byte, error) {
		return b, nil
	}))
	if err != nil {
		t.Fatalf("ParseRaw failed: %v", err)
	}
	m1.delims = delims

	j, err := m1.JSON(true)
	if err != nil {
		t.Fatalf("JSON failed: %v", err)
	}

	m2, err := ParseJSON(j)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	val, err := m2.QueryValue("PR1.2.1.1")
	if err != nil || val.String() != "SUB1" {
		t.Errorf("expected PR1.2.1.1 to be 'SUB1', got %v (err: %v)", val, err)
	}

	val, err = m2.QueryValue("PR1.2.1.2")
	if err != nil || val.String() != "SUB2" {
		t.Errorf("expected PR1.2.1.2 to be 'SUB2', got %v (err: %v)", val, err)
	}
}
