package money_accounts

import (
	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/utility"
)

func GenerateMoneyAccount() MoneyAccount {
	fields := GenerateAccountFields()
	balance := GenerateAccountBalace()
	moneyAccount := MoneyAccount{
		ID:                  uuid.New(),
		MoneyAccountFields:  fields,
		MoneyAccountBalance: balance,
	}
	return moneyAccount
}

func GenerateAccountFields() MoneyAccountFields {
	fields := MoneyAccountFields{
		Name:     utility.GetRandomString(25),
		IsCash:   utility.GetRandomBoolean(),
		Currency: utility.GetRandomCurrency(),
	}
	return fields
}

func GenerateAccountBalace() MoneyAccountBalance {
	balance := MoneyAccountBalance{
		Balance: utility.GetRandomBalance(),
	}
	return balance
}

func generateBadAccountFields() badAccountFields {
	badFields := badAccountFields{
		Name:     utility.GetRandomBoolean(),
		IsCash:   utility.GetRandomCurrency(),
		Currency: utility.GetRandomBoolean(),
	}
	return badFields
}
