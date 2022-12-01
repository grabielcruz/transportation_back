package money_accounts

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
)

// TestMoneyAccountServices contains a group of test related
// to the crud of moneyAccount
func TestMoneyAccountServices(t *testing.T) {
	envPath := filepath.Clean("../.env_test")
	sqlPath := filepath.Clean("../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
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

	deleteAllMoneyAccounts()

	t.Run("Create two money accounts and get and slice of accounts", func(t *testing.T) {
		CreateMoneyAccount(GenerateMoneyAccount())
		CreateMoneyAccount(GenerateMoneyAccount())
		moneyAccounts := GetMoneyAccounts()
		length := len(moneyAccounts)
		if length != 2 {
			t.Fatalf(`len(moneyAccounts) = %v, expected %v`, length, 2)
		}
	})

	deleteAllMoneyAccounts()

	t.Run("Create one money account and get it", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateMoneyAccount())
		obtainedMoneyAccount, err := GetOneMoneyAccount(createdMoneyAccount.ID)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if createdMoneyAccount.ID != obtainedMoneyAccount.ID {
			t.Fatalf(`GetOneMoneyAccount dit not returned the requested account, wanted %v, received %v`,
				createdMoneyAccount.ID, obtainedMoneyAccount.ID)
		}
	})

	deleteAllMoneyAccounts()

	t.Run("Error when getting unexisting account", func(t *testing.T) {
		var zeroUUID uuid.UUID
		_, err := GetOneMoneyAccount(zeroUUID)
		if err == nil {
			t.Fatalf(`Should generate error when getting money account with zero uuid`)
		}
	})

	t.Run("Create one money account and delete it", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateMoneyAccount())
		DeleteOneMoneyAccount(createdMoneyAccount.ID)
		_, err := GetOneMoneyAccount(createdMoneyAccount.ID)
		if err == nil {
			t.Fatalf(`Should receive error when requesting deleted money account`)
		}
	})

	deleteAllMoneyAccounts()
}
