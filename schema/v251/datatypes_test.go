package v251

import (
	"testing"

	"github.com/blushift-io/hl7v2/gen"
	"github.com/davecgh/go-spew/spew"
)

func TestGen(t *testing.T) {
	v := XAD{}
	if err := gen.Generate(&v); err != nil {
		t.Fatal(err)
	}

	spew.Dump(v)
}
