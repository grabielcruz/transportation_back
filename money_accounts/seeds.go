package money_accounts

import (
	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/utility"
)

func GenerateMoneyAccount() MoneyAccount {
	fields := MoneyAccountFields{
		Name:     utility.GetRandomString(25),
		IsCash:   utility.GetRandomBoolean(),
		Currency: utility.GetRandomCurrency(),
	}
	moneyAccount := MoneyAccount{
		ID:                 uuid.New(),
		Balance:            utility.GetRandomBalance(),
		MoneyAccountFields: fields,
	}
	return moneyAccount
}

func GenereatAccountFields() MoneyAccountFields {
	fields := MoneyAccountFields{
		Name:     utility.GetRandomString(25),
		IsCash:   utility.GetRandomBoolean(),
		Currency: utility.GetRandomCurrency(),
	}
	return fields
}
