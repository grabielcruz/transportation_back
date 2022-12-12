package money_accounts

import (
	"time"

	"github.com/google/uuid"
)

type MoneyAccount struct {
	ID uuid.UUID `json:"id"`
	MoneyAccountFields
	MoneyAccountBalance
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MoneyAccountFields struct {
	Name     string `json:"name"`
	IsCash   bool   `json:"is_cash"`
	Currency string `json:"currency"`
}

type badAccountFields struct {
	Name     bool   `json:"name"`
	IsCash   string `json:"is_cash"`
	Currency bool   `json:"currency"`
}

type MoneyAccountBalance struct {
	Balance float64 `json:"balance"`
}
