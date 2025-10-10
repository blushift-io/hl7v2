package v21

import "github.com/blushift-io/hl7v2"

type ADT_A01 struct {
	MSH hl7v2.MessageHeader `hl7:"MSH"`
	EVN EVN                 `hl7:"EVN"`
	PID any
	NK1 any
	PV1 any
	DG1 *DG1 `hl7:"DG1"`
}
