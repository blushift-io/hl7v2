package v251

type AD struct {
	StreetAddress              string `json:"street_address" gen:"street_address"`
	OtherDesignation           string `json:"other_designation"`
	City                       string `json:"city" gen:"city"`
	StateOrProvince            string `json:"state_or_province" gen:"state"`
	ZipOrPostalCode            string `json:"zip_or_postal_code" gen:"zip"`
	Country                    string `json:"country" gen:"country"`
	AddressType                string `json:"address_type"`
	OtherGeographicDesignation string `json:"other_geographic_designation"`
}

type XAD struct {
	StreetAddress              string `json:"street_address" gen:"street_address" hl7:"1"`
	OtherDesignation           string `json:"other_designation" hl7:"2"`
	City                       string `json:"city" gen:"city" hl7:"3"`
	StateOrProvince            string `json:"state_or_province" gen:"state" hl7:"4"`
	ZipOrPostalCode            string `json:"zip_or_postal_code" gen:"zip" hl7:"5"`
	Country                    string `json:"country" gen:"country" hl7:"6"`
	AddressType                string `json:"address_type" hl7:"7"`
	OtherGeographicDesignation string `json:"other_geographic_designation" hl7:"8"`
	CountyParishCode           string `json:"county_parish_code" hl7:"9"`
	CensusTract                string `json:"census_tract" hl7:"10"`
	AddressRepresentationCode  string `json:"address_representation_code" hl7:"11"`
	AddressValidityRange       string `json:"address_validity_range" hl7:"12"`
	EffectiveDate              string `json:"effective_date" hl7:"13"`
	ExpirationDate             string `json:"expiration_date" hl7:"14"`
}

type FN struct {
	Surname                        string `json:"surname" gen:"last_name" hl7:"1"`
	OwnSurnamePrefix               string `json:"own_surname_prefix" hl7:"2"`
	OwnSurname                     string `json:"own_surname" hl7:"3"`
	SurnamePrefixFromPartnerSpouse string `json:"surname_prefix_from_partner_spouse" hl7:"4"`
	SurnameFromPartnerSpouse       string `json:"surname_from_partner_spouse" hl7:"5"`
}

type XPN struct {
	FamilyName                                  FN     `json:"family_name"  hl7:"1"`
	GivenName                                   string `json:"given_name" gen:"first_name" hl7:"2"`
	SecondAndFurtherGivenNamesOrInitialsThereof string `json:"second_and_further_given_names_or_initials_thereof" hl7:"3"`
	Suffix                                      string `json:"suffix" hl7:"4"`
	Prefix                                      string `json:"prefix" hl7:"5"`
	Degree                                      string `json:"degree" hl7:"6"`
	NameTypeCode                                string `json:"name_type_code" hl7:"7"`
	NameRepresentationCode                      string `json:"name_representation_code" hl7:"8"`
	NameContext                                 string `json:"name_context" hl7:"9"`
	NameValidityRange                           string `json:"name_validity_range" hl7:"10"`
	NameAssemblyOrder                           string `json:"name_assembly_order" hl7:"11"`
	EffectiveDate                               string `json:"effective_date" hl7:"12"`
	ExpirationDate                              string `json:"expiration_date" hl7:"13"`
	ProfessionalSuffix                          string `json:"professional_suffix" hl7:"14"`
}
