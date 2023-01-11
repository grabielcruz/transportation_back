package money_accounts

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/utility"
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
		createdMoneyAccount, err := CreateMoneyAccount(accountFields)
		assert.Nil(t, err)
		assert.Equal(t, accountFields.Name, createdMoneyAccount.Name)
		assert.Equal(t, accountFields.Details, createdMoneyAccount.Details)
		assert.Equal(t, accountFields.Currency, createdMoneyAccount.Currency)
		assert.Equal(t, createdMoneyAccount.Balance, float64(0))
	})

	DeleteAllMoneyAccounts()

	t.Run("Create two money accounts and get an slice of accounts", func(t *testing.T) {
		CreateMoneyAccount(GenerateAccountFields())
		CreateMoneyAccount(GenerateAccountFields())
		moneyAccounts := GetMoneyAccounts()
		assert.Len(t, moneyAccounts, 2)
	})

	DeleteAllMoneyAccounts()

	t.Run("Create one money account and get it", func(t *testing.T) {
		createdMoneyAccount, err := CreateMoneyAccount(GenerateAccountFields())
		assert.Nil(t, err)
		obtainedMoneyAccount, err := GetOneMoneyAccount(createdMoneyAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdMoneyAccount.ID, obtainedMoneyAccount.ID)
	})

	DeleteAllMoneyAccounts()

	t.Run("Error when getting unexisting account", func(t *testing.T) {
		// with zero uuid
		zeroUUID := uuid.UUID{}
		_, err := GetOneMoneyAccount(zeroUUID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())

		// with random uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = GetOneMoneyAccount(randId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Create one money account and delete it", func(t *testing.T) {
		createdMoneyAccount, err := CreateMoneyAccount(GenerateAccountFields())
		assert.Nil(t, err)
		deletedId, err := DeleteOneMoneyAccount(createdMoneyAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdMoneyAccount.ID, deletedId.ID)
		_, err = GetOneMoneyAccount(createdMoneyAccount.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	DeleteAllMoneyAccounts()

	t.Run("Error when attempting to delete an unexisting account", func(t *testing.T) {
		// with zero uuid
		zeroUUID := uuid.UUID{}
		_, err := DeleteOneMoneyAccount(zeroUUID)
		assert.NotNil(t, err)

		// with random uuid
		randomUUID, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = DeleteOneMoneyAccount(randomUUID)
		assert.NotNil(t, err)
	})

	t.Run("It should create and update one money account", func(t *testing.T) {
		createFields := GenerateAccountFields()
		updateFields := GenerateAccountFields()
		createdAccount, err := CreateMoneyAccount(createFields)
		assert.Nil(t, err)
		updatedAccount, err := UpdateMoneyAccount(createdAccount.ID, updateFields)
		assert.Nil(t, err)
		assert.Equal(t, updatedAccount.ID, createdAccount.ID)
		assert.Equal(t, updateFields.Name, updatedAccount.Name)
		assert.Equal(t, updateFields.Currency, updatedAccount.Currency)
		assert.Equal(t, updateFields.Details, updatedAccount.Details)
		// assert.Greater(t, updatedAccount.UpdatedAt.Nanosecond(), updatedAccount.CreatedAt.Nanosecond())
	})

	DeleteAllMoneyAccounts()

	t.Run("It should generate error when trying to update an unexisting account", func(t *testing.T) {
		// with zero uuid
		zeroUUID := uuid.UUID{}
		zeroFields := MoneyAccountFields{}
		_, err := UpdateMoneyAccount(zeroUUID, zeroFields)
		assert.NotNil(t, err)

		// with random uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = UpdateMoneyAccount(randId, zeroFields)
		assert.NotNil(t, err)
	})

	t.Run("Create one money account and get its name", func(t *testing.T) {
		newMoneyAccount, err := CreateMoneyAccount(GenerateAccountFields())
		assert.Nil(t, err)
		name, err := getAccountsName(newMoneyAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, newMoneyAccount.Name, name)
	})

	t.Run("Error when getting unexisting money accounts name", func(t *testing.T) {
		// with zero uuid
		name, err := getAccountsName(uuid.UUID{})
		assert.Equal(t, "", name)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())

		// with random uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		name, err = getAccountsName(randId)
		assert.Equal(t, "", name)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Set accounts balance", func(t *testing.T) {
		newMoneyAccount, err := CreateMoneyAccount(GenerateAccountFields())
		assert.Nil(t, err)
		newBalance := utility.GetRandomBalance()
		updatedId, err := setAccountsBalance(newMoneyAccount.ID, newBalance)
		assert.Nil(t, err)
		assert.Equal(t, newMoneyAccount.ID, updatedId.ID)
		updatedAccount, err := GetOneMoneyAccount(newMoneyAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, newBalance, updatedAccount.Balance)
	})

	DeleteAllMoneyAccounts()

	t.Run("Error when updating unexisting account's balance", func(t *testing.T) {
		// with zero uuid
		zeroID := uuid.UUID{}
		newBalance := utility.GetRandomBalance()
		_, err := setAccountsBalance(zeroID, newBalance)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())

		// with random uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = setAccountsBalance(randId, newBalance)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Reset accounts balance", func(t *testing.T) {
		newMoneyAccount, err := CreateMoneyAccount(GenerateAccountFields())
		assert.Nil(t, err)
		updatedId, err := ResetAccountsBalance(newMoneyAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, newMoneyAccount.ID, updatedId.ID)
		updatedAccount, err := GetOneMoneyAccount(newMoneyAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, float64(0), updatedAccount.Balance)
	})

	DeleteAllMoneyAccounts()

	t.Run("Error when reseting unexisting account's balance", func(t *testing.T) {
		// with zero uuid
		zeroID := uuid.UUID{}
		_, err := ResetAccountsBalance(zeroID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())

		// with random uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = ResetAccountsBalance(randId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

}
