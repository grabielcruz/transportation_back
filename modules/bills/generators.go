package bills

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/utility"
)

func GenerateBillFields(person_id uuid.UUID) BillFields {
	fields := BillFields{
		PersonId:    person_id,
		Date:        time.Now(),
		Description: utility.GetRandomString(55),
		Currency:    utility.GetRandomCurrency(),
		Amount:      utility.GetRandomBalance(),
	}
	return fields
}

func GenerateBillToPayFields(person_id uuid.UUID) BillFields {
	fields := BillFields{
		PersonId:    person_id,
		Date:        time.Now(),
		Description: utility.GetRandomString(55),
		Currency:    utility.GetRandomCurrency(),
		Amount:      utility.GetRandomNegativeBalance(),
	}
	return fields
}

func GenerateBillToChargeFields(person_id uuid.UUID) BillFields {
	fields := BillFields{
		PersonId:    person_id,
		Date:        time.Now(),
		Description: utility.GetRandomString(55),
		Currency:    utility.GetRandomCurrency(),
		Amount:      utility.GetRandomPositiveBalance(),
	}
	return fields
}
