package transactions

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
)

type Transaction struct {
	ID uuid.UUID `json:"id"`
	TransactionFields
	PersonName string  `json:"person_name"`
	Balance    float64 `json:"balance"`
	common.Timestamps
}

type TransactionFields struct {
	AccountId   uuid.UUID `json:"account_id"`
	PersonId    uuid.UUID `json:"person_id"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
}

type TransationResponse struct {
	Transactions []Transaction `json:"transactions"`
	Pagination   common.Pagination
}

type badTransactionFields struct {
	AccountId   uuid.UUID `json:"account_id"`
	PersonId    uuid.UUID `json:"person_id"`
	Date        string    `json:"date"`
	Amount      string    `json:"amount"`
	Description bool      `json:"description"`
}

type badTransactionFieldsWithBadIds struct {
	AccountId   string `json:"account_id"`
	PersonId    string `json:"person_id"`
	Date        string `json:"date"`
	Amount      string `json:"amount"`
	Description bool   `json:"description"`
}
