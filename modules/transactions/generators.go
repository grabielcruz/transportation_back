package transactions

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/utility"
)

func GenerateTransactionFields(account_id uuid.UUID, person_id uuid.UUID) TransactionFields {
	fields := TransactionFields{
		AccountId:   account_id,
		PersonId:    person_id,
		Date:        time.Now(),
		Amount:      utility.GetRandomBalance(),
		Description: utility.GetRandomString(55),
	}
	return fields
}
