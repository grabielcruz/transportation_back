package money_accounts

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	"github.com/stretchr/testify/assert"
)

// TestMoneyAccountServices contains a group of test related
// to the crud of moneyAccount
func TestMoneyAccountServices(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()

	t.Run("Get empty slice of accounts initially", func(t *testing.T) {
		moneyAccounts := GetMoneyAccounts()
		assert.Len(t, moneyAccounts, 0)
	})

	t.Run("Create one money account", func(t *testing.T) {
		accountFields := GenerateAccountFields()
		createdMoneyAccount := CreateMoneyAccount(accountFields)
		assert.Equal(t, accountFields.Name, createdMoneyAccount.Name)
		assert.Equal(t, accountFields.Details, createdMoneyAccount.Details)
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
		zeroUUID := uuid.UUID{}
		_, err := GetOneMoneyAccount(zeroUUID)
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
	})

	t.Run("Create one money account and delete it", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		deletedId, err := DeleteOneMoneyAccount(createdMoneyAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdMoneyAccount.ID, deletedId.ID)
		_, err = GetOneMoneyAccount(createdMoneyAccount.ID)
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
	})

	deleteAllMoneyAccounts()

	t.Run("Error when attempting to delete an unexisting account", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		_, err := DeleteOneMoneyAccount(zeroUUID)
		assert.NotNil(t, err)
	})

	t.Run("It should create and update one money account", func(t *testing.T) {
		createFields := GenerateAccountFields()
		updateFields := GenerateAccountFields()
		createdAccount := CreateMoneyAccount(createFields)
		updatedAccount, err := UpdateMoneyAccount(createdAccount.ID, updateFields)
		assert.Nil(t, err)
		assert.Equal(t, updatedAccount.ID, createdAccount.ID)
		assert.Equal(t, updateFields.Name, updatedAccount.Name)
		assert.Equal(t, updateFields.Currency, updatedAccount.Currency)
		assert.Equal(t, updateFields.Details, updatedAccount.Details)
		assert.NotEqual(t, updatedAccount.CreatedAt.Nanosecond(), updatedAccount.UpdatedAt.Nanosecond())
	})

	deleteAllMoneyAccounts()

	t.Run("It should generate error when trying to update an unexisting account", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		zeroFields := MoneyAccountFields{}
		_, err := UpdateMoneyAccount(zeroUUID, zeroFields)
		assert.NotNil(t, err)
	})

	t.Run("It should create an account and update its balance", func(t *testing.T) {
		balance := GenerateAccountBalace()
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		updatedAccount, _ := updateMoneyAccountBalance(createdMoneyAccount.ID, balance)
		assert.Equal(t, updatedAccount.ID, createdMoneyAccount.ID)
		assert.Equal(t, balance, updatedAccount.Balance)
	})

	deleteAllMoneyAccounts()

	t.Run("Update balance time should be greater than creation time", func(t *testing.T) {
		balance := GenerateAccountBalace()
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		updatedAccount, _ := updateMoneyAccountBalance(createdMoneyAccount.ID, balance)
		assert.Equal(t, updatedAccount.ID, createdMoneyAccount.ID)
	})

	deleteAllMoneyAccounts()

	t.Run("It should get error when updating balance's account with wrong id", func(t *testing.T) {
		var zeroUUID uuid.UUID
		balance := GenerateAccountBalace()
		_, err := updateMoneyAccountBalance(zeroUUID, balance)
		assert.NotNil(t, err)
	})

	t.Run("It should modify account's balance and get it modified", func(t *testing.T) {
		balance := GenerateAccountBalace()
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		_, err := updateMoneyAccountBalance(createdMoneyAccount.ID, balance)
		assert.Nil(t, err)
		updatedBalance, err := GetAccountsBalance((createdMoneyAccount.ID))
		assert.Nil(t, err)
		assert.Equal(t, balance, updatedBalance)
	})

	deleteAllMoneyAccounts()

	t.Run("Get error when getting unexisting account's balance", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		balance, err := GetAccountsBalance(zeroUUID)
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
		assert.Zero(t, balance)
	})

	t.Run("Create an account and add an amount to its balance", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		newAmount := GenerateAccountBalace()
		updated, err := AddToBalance(createdMoneyAccount.ID, newAmount)
		assert.Nil(t, err)
		assert.Equal(t, newAmount, updated.Balance)
	})

	deleteAllMoneyAccounts()

	t.Run("Error when adding an amount to an unexisting account", func(t *testing.T) {
		newAmount := GenerateAccountBalace()
		nameAndBalance, err := AddToBalance(uuid.UUID{}, newAmount)
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
		assert.Equal(t, float64(0), nameAndBalance.Balance)
	})

	t.Run("Error when generating a negative balance", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		var newAmount float64 = -100
		nameAndBalance, err := AddToBalance(createdMoneyAccount.ID, newAmount)
		assert.NotNil(t, err)
		assert.Equal(t, float64(0), nameAndBalance.Balance)
		assert.Equal(t, "New balance can't be a negative number", err.Error())
	})

	deleteAllMoneyAccounts()

	t.Run("Add and substrat to account getting it back to zero", func(t *testing.T) {
		createdMoneyAccount := CreateMoneyAccount(GenerateAccountFields())
		newAmount := GenerateAccountBalace()
		_, err := AddToBalance(createdMoneyAccount.ID, newAmount)
		assert.Nil(t, err)
		nameAndBalance, err := AddToBalance(createdMoneyAccount.ID, float64(-1)*newAmount)
		assert.Nil(t, err)
		assert.Equal(t, float64(0), nameAndBalance.Balance)
	})
}
