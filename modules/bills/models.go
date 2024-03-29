package bills

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
)

type Bill struct {
	ID         uuid.UUID `json:"id"`
	PersonName string    `json:"person_name"`
	Status     string    `json:"status"`
	BillFields
	// Only closed by one of these
	TransactionId       uuid.UUID `json:"transaction_id"`
	BillCrossId         uuid.UUID `json:"bill_cross_id"`
	RevertTransactionId uuid.UUID `json:"revert_transaction_id"`
	// Only after closed
	PostNotes string `json:"post_notes"`
	// timestamps
	common.Timestamps
}

type BillFields struct {
	PersonId    uuid.UUID `json:"person_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Currency    string    `json:"currency"`
	Amount      float64   `json:"amount"`
	// either generated by one of these or none of these
	ParentTransactionId uuid.UUID `json:"parent_transaction_id"`
	ParentBillCrossId   uuid.UUID `json:"parent_bill_cross_id"`
}

type BillResponse struct {
	Bills          []Bill    `json:"bills"`
	FilterPersonId uuid.UUID `json:"filter_person_id"`
	common.Pagination
}
