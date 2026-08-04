package mock

import (
	"fmt"
	"reflect"

	fake "github.com/brianvoe/gofakeit/v7"
)

const (
	StreetAddressTag = "street_address"
	CityTag          = "city"
	StateTag         = "state"
	ZipTag           = "zip"
	CountryTag       = "country"
	LastNameTag      = "last_name"
	FirstNameTag     = "first_name"
)

type genFunc func(v reflect.Value) error

var genFuncs = map[string]genFunc{
	StreetAddressTag: genStreetAddress,
	//CityTag:          genCity,
	//StateTag:         genState,
	//ZipTag:           genZip,
	//CountryTag:       genCountry,
	//LastNameTag:      genLastName,
	//FirstNameTag:     genFirstName,
}

func genStreetAddress(v reflect.Value) error {
	if v.Elem().Kind() != reflect.String {
		return fmt.Errorf("genStreetAddress: expected string, got %s", v.Kind())
	}

	v.SetString(fake.Street())

	return nil
}
