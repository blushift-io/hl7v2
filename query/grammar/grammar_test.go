package grammar

import (
	"reflect"
	"testing"
)

func TestParseGrammar(t *testing.T) {
	table := []struct {
		q    string
		want []Grammar
	}{
		{"MSH [OBX]",
			[]Grammar{
				{
					ID: "MSH",
				},
			}},
		{"OBX {OBR}", []Grammar{
			{
				ID:        "OBX",
				Repeating: true,
				Children: []Grammar{
					{
						ID: "OBX",
					},
				},
			}},
		},
	}

	for _, tt := range table {
		got, err := New(tt.q)
		if err != nil {
			t.Errorf("error parsing grammar %s: %s", tt.q, err)
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}
