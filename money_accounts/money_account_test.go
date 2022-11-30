package money_accounts

import (
	"testing"

	"github.com/grabielcruz/transportation_back/database"
)

// TestMoneyAccountServices contains a group of test related
// to the crud of moneyAccount
func TestMoneyAccountServices(t *testing.T) {
	database.SetupDB()
	defer database.CloseConnection()

	t.Run("Get empty slice of accounts initially", func(t *testing.T) {
		var moneyAccounts []MoneyAccount
		moneyAccounts = GetMoneyAccounts()
		length := len(moneyAccounts)
		if length != 0 {
			t.Fatalf(`len(GetMoneyAccounts()) = %v, want 0`, length)
		}
	})

	t.Run("Create one money account", func(t *testing.T) {
		moneyAccount := GenerateMoneyAccount()
		createdMoneyAccount := CreateMoneyAccount(moneyAccount)
		if moneyAccount.Name != createdMoneyAccount.Name {
			t.Fatalf(`createdMoneyAccount.Name = %v, expected %v`, createdMoneyAccount.Name, moneyAccount.Name)
		}
		if moneyAccount.Balance != createdMoneyAccount.Balance {
			t.Fatalf(`createdMoneyAccount.Balance = %v, expected %v`, createdMoneyAccount.Balance, moneyAccount.Balance)
		}
		if moneyAccount.IsCash != createdMoneyAccount.IsCash {
			t.Fatalf(`createdMoneyAccount.IsCash = %v, expected %v`, createdMoneyAccount.IsCash, moneyAccount.IsCash)
		}
		if moneyAccount.Currency != createdMoneyAccount.Currency {
			t.Fatalf(`createdMoneyAccount.Currency = %v, expected %v`, createdMoneyAccount.Currency, moneyAccount.Currency)
		}
	})
}
