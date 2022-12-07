package money_accounts

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/stretchr/testify/assert"
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
		assert.Len(t, moneyAccounts, 0)
	})

	t.Run("Create one money account", func(t *testing.T) {
		accountFields := GenerateAccountFields()
		createdMoneyAccount := CreateMoneyAccount(accountFields)
		assert.Equal(t, accountFields.Name, createdMoneyAccount.Name)
		assert.Equal(t, accountFields.IsCash, createdMoneyAccount.IsCash)
		assert.Equal(t, accountFields.Currency, createdMoneyAccount.Currency)
		assert.Equal(t, createdMoneyAccount.Balance, float64(0))
	})

	deleteAllMoneyAccounts()

	t.Run("Create two money accounts and get an slice of accounts", func(t *testing.T) {
		CreateMoneyAccount(GenerateAccountFields())
		CreateMoneyAccount(GenerateAccountFields())
		moneyAccounts := GetMoneyAccounts()
		assert.Len(t, moneyAccounts, 2)
	})

	deleteAllMoneyAccounts()

	t.Run("Create one money account and get it", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		obtainedMoneyAccount, err := GetOneMoneyAccount(createdMoneyAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdMoneyAccount.ID, obtainedMoneyAccount.ID)
	})

	deleteAllMoneyAccounts()

	t.Run("Error when getting unexisting account", func(t *testing.T) {
		var zeroUUID uuid.UUID
		_, err := GetOneMoneyAccount(zeroUUID)
		assert.NotNil(t, err)
	})

	t.Run("Create one money account and delete it", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		DeleteOneMoneyAccount(createdMoneyAccount.ID)
		_, err := GetOneMoneyAccount(createdMoneyAccount.ID)
		assert.NotNil(t, err)
	})

	deleteAllMoneyAccounts()

	t.Run("Error when attempting to delete an unexisting account", func(t *testing.T) {
		var zeroUUID uuid.UUID
		_, err := DeleteOneMoneyAccount(zeroUUID)
		assert.NotNil(t, err)
	})

	t.Run("It should create and update one money account", func(t *testing.T) {
		createFields := GenerateAccountFields()
		updateFields := GenerateAccountFields()
		createdAccount := CreateMoneyAccount(createFields)
		updatedAccount, err := UpdateMoneyAccount(createdAccount.ID, updateFields)
		errors_handler.CheckError(err)
		assert.Equal(t, updatedAccount.ID, createdAccount.ID)
		assert.Equal(t, updateFields.Name, updatedAccount.Name)
		assert.Equal(t, updateFields.Currency, updatedAccount.Currency)
		assert.Equal(t, updateFields.IsCash, updatedAccount.IsCash)
	})

	deleteAllMoneyAccounts()

	t.Run("It should generate error when trying to update an unexisting account", func(t *testing.T) {
		var zeroUUID uuid.UUID
		var zeroFields MoneyAccountFields
		_, err := UpdateMoneyAccount(zeroUUID, zeroFields)
		assert.NotNil(t, err)
	})

	t.Run("It should create an account and update its balance", func(t *testing.T) {
		balance := GenerateAccountBalace()
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		updatedAccount, _ := UpdatedMoneyAccountsBalance(createdMoneyAccount.ID, balance)
		assert.Equal(t, createdMoneyAccount.ID, updatedAccount.ID)
		assert.Equal(t, updatedAccount.Balance, balance.Balance)
	})

	deleteAllMoneyAccounts()

	t.Run("It should get error when updating balance's account with wrong id", func(t *testing.T) {
		var zeroUUID uuid.UUID
		balance := GenerateAccountBalace()
		_, err := UpdatedMoneyAccountsBalance(zeroUUID, balance)
		assert.NotNil(t, err)
	})
}
