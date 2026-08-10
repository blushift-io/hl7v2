package builder_test

import (
	"testing"
	"time"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/builder"
)

func TestHeaderBuilderRefactored(t *testing.T) {
	mt := hl7v2.MessageType{Code: "ADT", Event: "A01", Structure: "ADT_A01"}
	testDate := time.Date(2026, 8, 10, 15, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		opts     []builder.HeaderOption
		validate func(t *testing.T, raw hl7v2.RawSegment)
	}{
		{
			name: "SendingApp and SendingFacility",
			opts: []builder.HeaderOption{
				builder.SendingApp("EPIC"),
				builder.SendingFacility("HOSPITAL_A"),
			},
			validate: func(t *testing.T, raw hl7v2.RawSegment) {
				if string(raw[3][0][0][0]) != "EPIC" {
					t.Errorf("expected EPIC sending app, got %s", raw[3][0][0][0])
				}
				if string(raw[4][0][0][0]) != "HOSPITAL_A" {
					t.Errorf("expected HOSPITAL_A sending facility, got %s", raw[4][0][0][0])
				}
			},
		},
		{
			name: "ReceivingApp and ReceivingFacility",
			opts: []builder.HeaderOption{
				builder.ReceivingApp("CERNER"),
				builder.ReceivingFacility("CLINIC_B"),
			},
			validate: func(t *testing.T, raw hl7v2.RawSegment) {
				if string(raw[5][0][0][0]) != "CERNER" {
					t.Errorf("expected CERNER receiving app, got %s", raw[5][0][0][0])
				}
				if string(raw[6][0][0][0]) != "CLINIC_B" {
					t.Errorf("expected CLINIC_B receiving facility, got %s", raw[6][0][0][0])
				}
			},
		},
		{
			name: "MessageDate, ControlID, ProcessingID",
			opts: []builder.HeaderOption{
				builder.MessageDate(testDate),
				builder.ControlID("MSG-12345"),
				builder.ProcessingID("P"),
			},
			validate: func(t *testing.T, raw hl7v2.RawSegment) {
				expectedDate := testDate.Format("20060102150405")
				if string(raw[7][0][0][0]) != expectedDate {
					t.Errorf("expected date %s, got %s", expectedDate, raw[7][0][0][0])
				}
				if string(raw[10][0][0][0]) != "MSG-12345" {
					t.Errorf("expected MSG-12345 control ID, got %s", raw[10][0][0][0])
				}
				if string(raw[11][0][0][0]) != "P" {
					t.Errorf("expected P processing ID, got %s", raw[11][0][0][0])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hb := builder.BuildHeader(mt, hl7v2.Version251, tt.opts...)
			raw := hb.Build()
			if len(raw) < 12 {
				t.Fatalf("expected header segment fields, got %d", len(raw))
			}
			tt.validate(t, raw)
		})
	}
}
