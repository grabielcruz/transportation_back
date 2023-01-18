package persons

import (
	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
)

type Person struct {
	ID uuid.UUID `json:"id"`
	PersonFields
	common.Timestamps
}

type PersonFields struct {
	Name     string `json:"name"`
	Document string `json:"document"`
}

type badPersonFields struct {
	Name     bool `json:"name"`
	Document bool `json:"document"`
}
