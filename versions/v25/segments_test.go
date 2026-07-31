package v25_test

import (
	"testing"

	"github.com/blushift-io/hl7v2/versions/v25"
)

func TestSegmentConstructorsAndSetters(t *testing.T) {
	pid := v25.NewPID().
		SetSetIDPID(v25.NewSI("1")).
		SetBirthPlace(v25.NewST("New York")).
		AddPatientName(v25.XPN{FamilyName: v25.FN{Surname: "Doe"}, GivenName: "John"}).
		AddPatientName(v25.XPN{FamilyName: v25.FN{Surname: "Doe"}, GivenName: "Johnny"})

	if string(*pid.SetIDPID) != "1" {
		t.Errorf("got SetIDPID %s, want 1", *pid.SetIDPID)
	}

	if string(*pid.BirthPlace) != "New York" {
		t.Errorf("got BirthPlace %s, want New York", *pid.BirthPlace)
	}

	if len(pid.PatientName) != 2 {
		t.Fatalf("got %d PatientNames, want 2", len(pid.PatientName))
	}

	if string(pid.PatientName[0].GivenName) != "John" {
		t.Errorf("got %s, want John", pid.PatientName[0].GivenName)
	}

	if string(pid.PatientName[1].GivenName) != "Johnny" {
		t.Errorf("got %s, want Johnny", pid.PatientName[1].GivenName)
	}
}

func TestMSHConstructor(t *testing.T) {
	msgType := v25.MSG{
		MessageCode: "ADT",
		TriggerEvent: "A01",
		MessageStructure: "ADT_A01",
	}

	msh := v25.NewMSH(&msgType)

	if string(*msh.FieldSeparator) != "|" {
		t.Errorf("got FieldSeparator %s, want |", *msh.FieldSeparator)
	}

	if string(*msh.EncodingCharacters) != "^~\\&" {
		t.Errorf("got EncodingCharacters %s, want ^~\\&", *msh.EncodingCharacters)
	}

	if string(msh.DateTimeOfMessage.Time) == "" {
		t.Errorf("expected non-empty DateTimeOfMessage")
	}

	if string(*msh.MessageControlID) == "" {
		t.Errorf("expected non-empty MessageControlID")
	}

	if string(msh.VersionID.VersionId) != "2.5" {
		t.Errorf("got VersionID %s, want 2.5", msh.VersionID.VersionId)
	}

	if string(msh.MessageType.MessageCode) != "ADT" {
		t.Errorf("got MessageType.MessageCode %s, want ADT", msh.MessageType.MessageCode)
	}
}
