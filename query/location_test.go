package query

import "testing"

func TestParseLocation(t *testing.T) {
	table := []struct {
		q    string
		want Location
	}{
		{"MSH.1",
			Location{
				Segment: "MSH",
				Field:   1,
			}},
		{"OBX[1].1.1", Location{
			Segment:    "OBX",
			SegmentRep: 1,
			Field:      1,
			Component:  1,
		}},
	}

	for _, tt := range table {
		got, err := ParseLocation(tt.q)
		if err != nil {
			t.Errorf("error parsing location %s: %s", tt.q, err)
		}

		if got != tt.want {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}
