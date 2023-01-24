package person_accounts

import (
	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
)

type PersonAccount struct {
	ID       uuid.UUID `json:"id"`
	PersonId uuid.UUID `json:"person_id"`
	PersonAccountFields
	common.Timestamps
}

type PersonAccountFields struct {
	UpdatePersonAccountFields
	Currency string `json:"currency"`
}

type UpdatePersonAccountFields struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type badPersonAccountFields struct {
	Name        bool `json:"name"`
	Description bool `json:"description"`
	Currency    bool `json:"currency"`
}
