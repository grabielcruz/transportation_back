package money_accounts

import (
	"github.com/grabielcruz/transportation_back/utility"
)

func GenerateAccountFields() MoneyAccountFields {
	fields := MoneyAccountFields{
		Name:     utility.GetRandomString(25),
		Details:  utility.GetRandomString(45),
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
		Details:  utility.GetRandomBoolean(),
		Currency: utility.GetRandomBoolean(),
	}
	return badFields
}
