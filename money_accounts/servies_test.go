package money_accounts

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
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
		accountFields := GenerateAccountFields()
		createdMoneyAccount := CreateMoneyAccount(accountFields)
		if accountFields.Name != createdMoneyAccount.Name {
			t.Fatalf(`createdMoneyAccount.Name = %v, expected %v`, createdMoneyAccount.Name, accountFields.Name)
		}
		if accountFields.IsCash != createdMoneyAccount.IsCash {
			t.Fatalf(`createdMoneyAccount.IsCash = %v, expected %v`, createdMoneyAccount.IsCash, accountFields.IsCash)
		}
		if accountFields.Currency != createdMoneyAccount.Currency {
			t.Fatalf(`createdMoneyAccount.Currency = %v, expected %v`, createdMoneyAccount.Currency, accountFields.Currency)
		}
		if createdMoneyAccount.Balance != 0 {
			t.Fatalf(`createdMoneyAccount.Balance = %v, expected %v`, createdMoneyAccount.Currency, 0)
		}
	})

	deleteAllMoneyAccounts()

	t.Run("Create two money accounts and get and slice of accounts", func(t *testing.T) {
		CreateMoneyAccount(GenerateAccountFields())
		CreateMoneyAccount(GenerateAccountFields())
		moneyAccounts := GetMoneyAccounts()
		length := len(moneyAccounts)
		if length != 2 {
			t.Fatalf(`len(moneyAccounts) = %v, expected %v`, length, 2)
		}
	})

	deleteAllMoneyAccounts()

	t.Run("Create one money account and get it", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
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
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		DeleteOneMoneyAccount(createdMoneyAccount.ID)
		_, err := GetOneMoneyAccount(createdMoneyAccount.ID)
		if err == nil {
			t.Fatalf(`Should receive error when requesting deleted money account`)
		}
	})

	deleteAllMoneyAccounts()

	t.Run("Error when attempting to delete an unexisting account", func(t *testing.T) {
		var zeroUUID uuid.UUID
		_, err := DeleteOneMoneyAccount(zeroUUID)
		if err == nil {
			t.Fatalf(`Should get an error when attempting to delete an unexisting account`)
		}
	})

	t.Run("It should create and update one money account", func(t *testing.T) {
		createFields := GenerateAccountFields()
		updateFields := GenerateAccountFields()
		createdAccount := CreateMoneyAccount(createFields)
		updatedAccount, err := UpdateMoneyAccount(createdAccount.ID, updateFields)
		errors_handler.CheckError(err)
		if updatedAccount.ID != createdAccount.ID {
			t.Fatalf(`UpdateMoneyAccount did not return same account's id, wanted %v, got %v`, createdAccount.ID, updatedAccount.ID)
		}
		if updateFields.Name != updatedAccount.Name {
			t.Fatalf(`UpdatedMoneyAccount did not updated account's name, wanted %v, got %v`, updateFields.Name, updatedAccount.Name)
		}
		if updateFields.Currency != updatedAccount.Currency {
			t.Fatalf(`UpdatedMoneyAccount did not updated account's currency, wanted %v, got %v`, updateFields.Currency, updatedAccount.Currency)
		}
		if updateFields.IsCash != updatedAccount.IsCash {
			t.Fatalf(`UpdatedMoneyAccount did not updated account's IsCash property, wanted %v, got %v`, updateFields.IsCash, updatedAccount.IsCash)
		}
	})

	deleteAllMoneyAccounts()

	t.Run("It should generate error when trying to update an unexisting account", func(t *testing.T) {
		var zeroUUID uuid.UUID
		var zeroFields MoneyAccountFields
		_, err := UpdateMoneyAccount(zeroUUID, zeroFields)
		if err == nil {
			t.Fatalf(`UpdateMoneyAccount should generate error, instead generated nil`)
		}
	})

	t.Run("It should create an account and update its balance", func(t *testing.T) {
		balance := GenerateAccountBalace()
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		updatedAccount, _ := UpdatedMoneyAccountsBalance(createdMoneyAccount.ID, balance)
		if createdMoneyAccount.ID != updatedAccount.ID {
			t.Fatalf(`Updated account does not have the right id, want %v, got %v`, createdMoneyAccount.ID, updatedAccount.ID)
		}
		if updatedAccount.Balance != balance.Balance {
			t.Fatalf(`Updated account' balance is %v, expected %v`, updatedAccount.Balance, balance)
		}
	})

	deleteAllMoneyAccounts()

	t.Run("It should get error when updating balance's account with wrong id", func(t *testing.T) {
		var zeroUUID uuid.UUID
		balance := GenerateAccountBalace()
		_, err := UpdatedMoneyAccountsBalance(zeroUUID, balance)
		if err == nil {
			t.Fatalf(`UpdateMoneyAccountsBalance should generate error, instead generated nil`)
		}
	})
}
