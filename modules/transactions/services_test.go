package transactions

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/stretchr/testify/assert"
)

func TestTransactionServices(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	account := money_accounts.CreateMoneyAccount(money_accounts.GenerateAccountFields())
	person := persons.CreatePerson(persons.GeneratePersonFields())

	t.Run("Get transaction response with zero transactions", func(t *testing.T) {
		transactions := GetTransactions(account.ID, Limit, Offset)
		assert.Len(t, transactions.Transactions, 0)
		assert.Equal(t, transactions.Pagination.Count, 0)
		assert.Equal(t, transactions.Pagination.Offset, Offset)
		assert.Equal(t, transactions.Pagination.Limit, Limit)
	})

	t.Run("Create one transaction without person", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, uuid.UUID{})
		newTransaction, err := CreateTransaction(transactionFields)
		assert.Nil(t, err)
		updatedAccount, err := money_accounts.GetOneMoneyAccount(newTransaction.AccountId)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.Balance, updatedAccount.Balance)
		assert.Equal(t, newTransaction.AccountId, updatedAccount.ID)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction with a person", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		newTransaction, err := CreateTransaction(transactionFields)
		assert.Nil(t, err)
		updatedAccount, err := money_accounts.GetOneMoneyAccount(newTransaction.AccountId)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.Balance, updatedAccount.Balance)
		assert.Equal(t, newTransaction.AccountId, updatedAccount.ID)
		assert.Equal(t, newTransaction.PersonName, person.Name)
		assert.Equal(t, newTransaction.PersonId, person.ID)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when creating transaction with wrong account's id", func(t *testing.T) {
		zeroId := uuid.UUID{}
		transactionFields := GenerateTransactionFields(zeroId, zeroId)
		_, err := CreateTransaction(transactionFields)
		assert.NotNil(t, err)
		assert.Equal(t, "Could not get balance from account", err.Error())
	})

	t.Run("Error when generating negative balance", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		transactionFields.Amount *= -1
		_, err := CreateTransaction(transactionFields)
		assert.NotNil(t, err)
		assert.Equal(t, "Transaction should not generate a negative balance", err.Error())
		updatedAccount, err := money_accounts.GetOneMoneyAccount(transactionFields.AccountId)
		assert.Nil(t, err)
		// accounts balance should remain unmodified, which means it is equal to zero
		assert.Equal(t, float64(0), updatedAccount.Balance)
	})

	t.Run("Transaction should have a description", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		transactionFields.Description = ""
		_, err := CreateTransaction(transactionFields)
		assert.NotNil(t, err)
		assert.Equal(t, "Transaction should have a description", err.Error())
		updatedAccount, err := money_accounts.GetOneMoneyAccount(transactionFields.AccountId)
		assert.Nil(t, err)
		// accounts balance should remain unmodified, which means it is equal to zero
		assert.Equal(t, float64(0), updatedAccount.Balance)
	})

	t.Run("Transaction should have an amount greater than zero", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		transactionFields.Amount = float64(0)
		_, err := CreateTransaction(transactionFields)
		assert.NotNil(t, err)
		assert.Equal(t, "Amount should be greater than zero", err.Error())
		updatedAccount, err := money_accounts.GetOneMoneyAccount(transactionFields.AccountId)
		assert.Nil(t, err)
		assert.Equal(t, float64(0), updatedAccount.Balance)
	})

	t.Run("Create one transaction and get it", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		newTransaction, err := CreateTransaction(transactionFields)
		assert.Nil(t, err)
		transactions := GetTransactions(account.ID, Limit, Offset)
		assert.Equal(t, Offset, transactions.Pagination.Offset)
		assert.Equal(t, Limit, transactions.Pagination.Limit)
		assert.Equal(t, 1, transactions.Pagination.Count)
		assert.Equal(t, newTransaction, transactions.Transactions[0])
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Execute 1000 transactions and get accounts balance right", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(1000)
		sum := utility.GetSumOfAmounts(amounts)
		for i, v := range amounts {
			personId := person.ID
			if i%20 == 0 {
				personId = uuid.UUID{}
			}
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}
		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)
		assert.Equal(t, sum, updatedAccount.Balance)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	// at the end of all transactions services tests
	money_accounts.DeleteAllMoneyAccounts()
	persons.DeleteAllPersons()
}
