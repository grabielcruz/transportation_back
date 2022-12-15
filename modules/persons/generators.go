package persons

import "github.com/grabielcruz/transportation_back/utility"

func GeneratePersonFields() PersonFields {
	fields := PersonFields{
		Name:     utility.GetRandomString(20),
		Document: utility.GetRandomString(16),
	}
	return fields
}

func generateBadPersonFields() badPersonFields {
	fields := badPersonFields{
		Name:     utility.GetRandomBoolean(),
		Document: utility.GetRandomBoolean(),
	}
	return fields
}
