package hl7v2

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNewMessage(t *testing.T) {

	msg, err := NewMessageFromFile("./fixtures/ORU_R01_1.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(msg.header)

	walkElements(msg)

}

func walkElements(root Element) {
	it := NewIterator(root)
	for el := it.Next(); el != nil; el = it.Next() {
		if el.Type() < ElementComponent || len(el.Children()) > 1 {
			walkElements(el)
		} else {
			fmt.Printf("%s - %s\n", el.Location().String(), el.Value().String())
		}
	}
}

func TestMessageGrammar(t *testing.T) {
	msg, err := NewMessageFromFile("./fixtures/ORM_O01_1.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	els, err := msg.Select("{OBR DG1}")
	if err != nil {
		t.Fatal(err)
	}

	for _, seg := range els {
		walkElements(seg)
	}
}
