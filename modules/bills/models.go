package bills

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
)

type Bill struct {
	ID         uuid.UUID `json:"id"`
	PersonName string    `json:"person_name"`
	BillFields
	Pending float64 `json:"pending"`
	common.Timestamps
}

type BillFields struct {
	PersonId    uuid.UUID `json:"person_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Currency    string    `json:"currency"`
	Amount      float64   `json:"amount"`
}

type BillResponse struct {
	Bills          []Bill    `json:"bills"`
	FilterPersonId uuid.UUID `json:"filter_person_id"`
	common.Pagination
}
