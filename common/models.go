package common

import (
	"time"

	"github.com/google/uuid"
)

type ID struct {
	ID uuid.UUID `json:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
