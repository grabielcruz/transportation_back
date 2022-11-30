package money_accounts

import (
	"time"

	"github.com/google/uuid"
)

type MoneyAccount struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	IsCash    bool      `json:"is_cash"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
