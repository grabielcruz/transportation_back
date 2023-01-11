package currencies

import (
	"path/filepath"
	"testing"

	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/stretchr/testify/assert"
)

func TestCurrenciesServices(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()

	t.Run("Can get initially an array with two currencies", func(t *testing.T) {
		currencies := GetCurrencies()
		assert.Len(t, currencies, 2)
		assert.Equal(t, "VED", currencies[0])
		assert.Equal(t, "USD", currencies[1])
	})

	t.Run("Can create a currency", func(t *testing.T) {
		newCurrency := "ABC"
		createdCurrency, err := CreateCurrency(newCurrency)
		assert.Nil(t, err)
		assert.Equal(t, newCurrency, createdCurrency)
		currencies := GetCurrencies()
		assert.Len(t, currencies, 3)
		assert.Equal(t, "VED", currencies[0])
		assert.Equal(t, "USD", currencies[1])
		assert.Equal(t, newCurrency, currencies[2])
	})

	resetCurrencies()

	t.Run("Error when creating repeated currency", func(t *testing.T) {
		newCurrency := "VED"
		_, err := CreateCurrency(newCurrency)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.CU003, err.Error())
	})

	t.Run("Error when creating empty currency", func(t *testing.T) {
		newCurrency := ""
		_, err := CreateCurrency(newCurrency)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.CU002, err.Error())
	})

	t.Run("Can create a currency, then delete it", func(t *testing.T) {
		newCurrency := "ABC"
		createdCurrency, err := CreateCurrency(newCurrency)
		assert.Nil(t, err)
		assert.Equal(t, newCurrency, createdCurrency)

		// deleting
		deletedCurrency, err := DeleteCurrency(newCurrency)
		assert.Nil(t, err)
		assert.Equal(t, newCurrency, deletedCurrency)

		// checking
		currencies := GetCurrencies()
		assert.Len(t, currencies, 2)
		assert.Equal(t, "VED", currencies[0])
		assert.Equal(t, "USD", currencies[1])
	})

	resetCurrencies()

	t.Run("Error when deleting unexisting currency", func(t *testing.T) {
		currency := "KKK"
		_, err := DeleteCurrency(currency)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Error when trying to delete VED or USD currencies", func(t *testing.T) {
		_, err := DeleteCurrency("VED")
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.CU001, err.Error())
		_, err = DeleteCurrency("USD")
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.CU001, err.Error())
	})

	t.Run("Error when trying to delete currency associated with a money account", func(t *testing.T) {
		createdCurrency, err := CreateCurrency("ABC")
		assert.Nil(t, err)

		accountsFields := money_accounts.GenerateAccountFields()
		accountsFields.Currency = createdCurrency
		newMoneyAccount, err := money_accounts.CreateMoneyAccount(accountsFields)
		assert.Nil(t, err)
		assert.Equal(t, newMoneyAccount.Currency, createdCurrency)

		_, err = DeleteCurrency(createdCurrency)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.CU004, err.Error())
	})

	t.Run("Error when deleting zero currency", func(t *testing.T) {
		currency := "000"
		_, err := DeleteCurrency(currency)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.CU002, err.Error())
	})

	resetCurrencies()
	money_accounts.DeleteAllMoneyAccounts()
}
