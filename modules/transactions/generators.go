package transactions

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/utility"
)

// GenerateRandomTransactionFields generates wheter a Incoming transaction
// or a outgoing transaction. Fee number and amount sign diferentiates them
func GenerateRandomTransactionFields(account_id uuid.UUID) TransactionFields {
	if utility.GetRandomBoolean() {
		return GenerateIncomingTransactionFields(account_id)
	}
	return GenerateOutgoingTransactionFields(account_id)
}

func GenerateIncomingTransactionFields(account_id uuid.UUID) TransactionFields {
	fields := TransactionFields{
		AccountId:   account_id,
		Date:        time.Now(),
		Amount:      utility.GetRandomPositiveBalance(),
		Fee:         0,
		Description: utility.GetRandomString(55),
	}
	return fields
}

func GenerateOutgoingTransactionFields(account_id uuid.UUID) TransactionFields {
	fields := TransactionFields{
		AccountId:   account_id,
		Date:        time.Now(),
		Amount:      utility.GetRandomNegativeBalance(),
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
