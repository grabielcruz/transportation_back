package transactions

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/utility"
)

func GenerateTransactionFields(account_id uuid.UUID) TransactionFields {
	fields := TransactionFields{
		AccountId:   account_id,
		Date:        time.Now(),
		Amount:      utility.GetRandomBalance(),
		Fee:         utility.GetRandomFee(),
		Description: utility.GetRandomString(55),
	}
	return fields
}

func generateBadTransactionFields(account_id uuid.UUID) badTransactionFields {
	badFields := badTransactionFields{
		AccountId:   account_id,
		Date:        utility.GetRandomString(10),
		Amount:      utility.GetRandomString(10),
		Description: utility.GetRandomBoolean(),
	}
	return badFields
}

func generateBadTransactionFieldsWithBadIds(account_id uuid.UUID) badTransactionFieldsWithBadIds {
	badFields := badTransactionFieldsWithBadIds{
		AccountId:   "absce",
		Date:        utility.GetRandomString(10),
		Amount:      utility.GetRandomString(10),
		Description: utility.GetRandomBoolean(),
	}
	return badFields
}
