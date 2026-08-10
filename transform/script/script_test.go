package script

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blushift-io/hl7v2"
	"github.com/stretchr/testify/assert"
)

var testScript = `
function handleMessage(m) { 
	var msh = m.getSegment("MSH")
	if (!msh) {
		console.log("no msh")
	} else {
		var evt = msh.getComponent(9,1)
		if (evt == "ADT") {
			msh.setField(12, "2.1")
			msh.addField("ABCD")
		}
	}

	try {
		var rm = parseMessage(m.toString())
		
		var fv = rm.queryValue("MSH.9.1")
		console.log(fv.string())
		
	} catch (e) {
		console.log(e)
	}

	m.addSegment("ZZX", "1", "test", "foo", "bar");
	var res = m.toString();

	return res
}

var p = new parser();
var m = p.parse(rawMessage);

handleMessage(m);
`

var testMsg = `MSH|^~\&|ADT1|GOOD HEALTH HOSPITAL|GHH LAB, INC.|GOOD HEALTH HOSPITAL|198808181126|SECURITY|ADT^A01^ADT_A01|MSG00001|P|2.8||
EVN|A01|200708181123||
PID|1||PATID1234^5^M11^ADT1^MR^GOOD HEALTH HOSPITAL~123456789^^^USSSA^SS||EVERYMAN^ADAM^A^III||19610615|M||C|2222 HOME STREET^^GREENSBORO^NC^27401-1020|GL|(555) 555-2004|(555)555-2004||S||PATID12345001^2^M10^ADT1^AN^A|444333333|987654^NC|
NK1|1|NUCLEAR^NELDA^W|SPO^SPOUSE||||NK^NEXT OF KIN
PV1|1|I|2000^2012^01||||004777^ATTEND^AARON^A|||SUR||||ADM|A0|`

func TestNewTransform(t *testing.T) {
	tr, err := NewTransform(testScript)
	if err != nil {
		t.Fatal(err)
	}

	m := hl7v2.ReplaceLineEndings([]byte(testMsg))

	res, err := tr(m)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := hl7v2.NewMessage(res, hl7v2.FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	z, err := msg.Query("ZZX.2")
	if err != nil {
		t.Fatalf("output: \n%s\n error %s\n", res, err)
	}

	assert.Equal(t, "test", z.Value().String())
}

func TestNewScript(t *testing.T) {
	scr, err := newScript(testScript)
	if err != nil {
		t.Fatal(err)
	}

	for i, l := range strings.Split(scr, "\n") {
		fmt.Printf("%d \t%s\n", i, l)
	}
}
