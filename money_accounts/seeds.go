package money_accounts

import (
	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/utility"
)

func GenerateMoneyAccount() MoneyAccount {
	moneyAccount := MoneyAccount{
		ID:       uuid.New(),
		Name:     utility.GetRandomString(25),
		Balance:  utility.GetRandomBalance(),
		IsCash:   utility.GetRandomBoolean(),
		Currency: utility.GetRandomCurrency(),
	}
	return moneyAccount
}
