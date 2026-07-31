package v25_test

import (
	"testing"

	"github.com/blushift-io/hl7v2/versions/v25"
)

func TestDatatypeConstructorsAndSetters(t *testing.T) {
	ce := v25.NewCE().
		SetIdentifier("1001").
		SetText("Glucose").
		SetNameOfCodingSystem("LN").
		SetAlternateIdentifier("GLU").
		SetAlternateText("Glucose Test").
		SetNameOfAlternateCodingSystem("L")

	if string(ce.Identifier) != "1001" {
		t.Errorf("got %s, want 1001", ce.Identifier)
	}
	if string(ce.Text) != "Glucose" {
		t.Errorf("got %s, want Glucose", ce.Text)
	}
	if string(ce.NameOfCodingSystem) != "LN" {
		t.Errorf("got %s, want LN", ce.NameOfCodingSystem)
	}

	ad := v25.NewAD().
		SetStreetAddress("123 Main St").
		SetOtherDesignation("Apt 4B").
		SetCity("Metropolis").
		SetStateOrProvince("NY").
		SetZipOrPostalCode("10001").
		SetCountry("USA").
		SetAddressType("H").
		SetOtherGeographicDesignation("Sub")

	if string(ad.StreetAddress) != "123 Main St" {
		t.Errorf("got %s, want 123 Main St", ad.StreetAddress)
	}
	if string(ad.City) != "Metropolis" {
		t.Errorf("got %s, want Metropolis", ad.City)
	}
	if string(ad.ZipOrPostalCode) != "10001" {
		t.Errorf("got %s, want 10001", ad.ZipOrPostalCode)
	}
}
