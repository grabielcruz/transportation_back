package person_accounts

import "github.com/grabielcruz/transportation_back/utility"

func GeneratePersonAccountFields() PersonAccountFields {
	fields := PersonAccountFields{
		Currency: utility.GetRandomCurrency(),
	}
	fields.Name = utility.GetRandomString(10)
	fields.Description = utility.GetRandomString(25)
	return fields
}

func GenerateUpdatePersonAccountFields() UpdatePersonAccountFields {
	fields := UpdatePersonAccountFields{
		Name:        utility.GetRandomString(10),
		Description: utility.GetRandomString(25),
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
