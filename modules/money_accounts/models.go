package money_accounts

import (
	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
)

type MoneyAccount struct {
	ID uuid.UUID `json:"id"`
	MoneyAccountFields
	MoneyAccountBalance
	common.Timestamps
}

type MoneyAccountFields struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Details  string `json:"details"`
}

type badAccountFields struct {
	Name     bool `json:"name"`
	Details  bool `json:"details"`
	Currency bool `json:"currency"`
}

type MoneyAccountBalance struct {
	Balance float64 `json:"balance"`
}
