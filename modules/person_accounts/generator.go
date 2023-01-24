package person_accounts

import "github.com/grabielcruz/transportation_back/utility"

func GeneratePersonAccountFields() PersonAccountFields {
	fields := PersonAccountFields{
		Name:        utility.GetRandomString(10),
		Description: utility.GetRandomString(25),
		Currency:    utility.GetRandomCurrency(),
	}
	return fields
}

func generateBadPersonAccountFields() badPersonAccountFields {
	fields := badPersonAccountFields{
		Name:        utility.GetRandomBoolean(),
		Description: utility.GetRandomBoolean(),
		Currency:    utility.GetRandomBoolean(),
	}
	return fields
}
