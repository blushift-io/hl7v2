package v25_test

import (
	"testing"

	"github.com/blushift-io/hl7v2/versions/v25"
)

func TestValidation(t *testing.T) {
	// 1. MSH Segment initial validation should pass because NewMSH initializes all required fields
	msh := v25.NewMSH(v25.NewMSG().SetMessageCode("ADT").SetTriggerEvent("A01").SetMessageStructure("ADT_A01"))
	if err := msh.Validate(); err != nil {
		t.Fatalf("expected initialized MSH to be valid, got: %v", err)
	}

	// Unsetting a required field on MSH should cause Validate() to fail
	msh.FieldSeparator = nil
	if err := msh.Validate(); err == nil {
		t.Errorf("expected MSH to fail validation when FieldSeparator is nil")
	}

	// 2. ADT_A01 message validation
	msg := v25.NewADT_A01()
	// EVN and PID are required in ADT_A01 v2.5
	if err := msg.Validate(); err == nil {
		t.Errorf("expected unpopulated ADT_A01 to fail validation due to missing required segments")
	}

	// Populate required segments on ADT_A01
	evn := v25.NewEVN().SetEventTypeCode(v25.NewID("A01")).SetRecordedDateTime(v25.NewTS().SetTime("20260731150000"))
	pid := v25.NewPID().SetSetIDPID(v25.NewSI("1")).
		AddPatientIdentifierList(*v25.NewCX().SetIdNumber("12345")).
		AddPatientName(*v25.NewXPN().SetFamilyName(*v25.NewFN().SetSurname("Doe")).SetGivenName("John"))

	pv1 := v25.NewPV1().SetPatientClass(v25.Table0004Inpatient.IS())

	msg.SetEVN(evn).SetPID(pid).SetPV1(pv1)

	if err := msg.Validate(); err != nil {
		t.Fatalf("expected populated ADT_A01 to be valid, got: %v", err)
	}
}
