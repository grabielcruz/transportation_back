package transactions

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/config"
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
	account, err := money_accounts.CreateMoneyAccount(money_accounts.GenerateAccountFields())
	assert.Nil(t, err)
	person, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)

	t.Run("Get transaction response with zero transactions initially", func(t *testing.T) {
		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, transactions.Transactions, 0)
		assert.Equal(t, transactions.Count, 0)
		assert.Equal(t, transactions.Offset, config.Offset)
		assert.Equal(t, transactions.Limit, config.Limit)
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

	t.Run("Error when creating transaction with unexisting account", func(t *testing.T) {
		zeroId := uuid.UUID{}
		transactionFields := GenerateTransactionFields(zeroId, zeroId)
		_, err := CreateTransaction(transactionFields)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR001, err.Error())
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when generating negative balance", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		transactionFields.Amount *= -1
		_, err := CreateTransaction(transactionFields)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR002, err.Error())
		updatedAccount, err := money_accounts.GetOneMoneyAccount(transactionFields.AccountId)
		assert.Nil(t, err)
		// accounts balance should remain unmodified, which means it is equal to zero
		assert.Equal(t, float64(0), updatedAccount.Balance)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction and get it in paginated response", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		newTransaction, err := CreateTransaction(transactionFields)
		assert.Nil(t, err)
		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Equal(t, config.Offset, transactions.Offset)
		assert.Equal(t, config.Limit, transactions.Limit)
		assert.Equal(t, 1, transactions.Count)
		assert.Equal(t, newTransaction, transactions.Transactions[0])
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction and get it with single response", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		newTransaction, err := CreateTransaction(transactionFields)
		assert.Nil(t, err)
		transaction, err := GetTransaction(newTransaction.ID)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.ID, transaction.ID)
		assert.Equal(t, newTransaction.AccountId, transaction.AccountId)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, newTransaction.Balance, transaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt, transaction.CreatedAt)
		assert.Equal(t, newTransaction.UpdatedAt, transaction.UpdatedAt)
		assert.Equal(t, newTransaction.Date, transaction.Date)
		assert.Equal(t, newTransaction.Description, transaction.Description)
		assert.Equal(t, newTransaction.PersonId, transaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, transaction.PersonName)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when getting non registered transaction", func(t *testing.T) {
		_, err := GetTransaction(uuid.UUID{})
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Execute 100 transactions and get accounts balance right", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(100)
		sum := utility.GetSumOfAmounts(amounts)
		for i, v := range amounts {
			personId := person.ID
			if i%5 == 0 {
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

	t.Run("Execute 10 transaction and the first transaction in the slice should be the last on execution", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(10)
		for i, v := range amounts {
			personId := person.ID
			if i%2 == 0 {
				personId = uuid.UUID{}
			}
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}
		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)
		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Equal(t, transactions.Transactions[0].Balance, updatedAccount.Balance)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Execute 51 transaction and get in last page the initial transaction, and count equal 51", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(51)
		for i, v := range amounts {
			personId := person.ID
			if i%3 == 0 {
				personId = uuid.UUID{}
			}
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}
		transactions, err := GetTransactions(account.ID, config.Limit, 50)
		assert.Nil(t, err)
		assert.Equal(t, transactions.Transactions[0].Balance, transactions.Transactions[0].Amount)
		assert.Equal(t, 51, transactions.Count)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create 7 transactions, update the last one and check sequence", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(7)
		for i, v := range amounts {
			personId := person.ID
			if i%3 == 0 {
				personId = uuid.UUID{}
			}
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}
		updatedAccountFirst, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)

		transactionResponseFirst, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)

		assert.Equal(t, updatedAccountFirst.Balance, transactionResponseFirst.Transactions[0].Balance)
		generatedBalanceFromPreviousTransaction := utility.RoundToTwoDecimalPlaces(transactionResponseFirst.Transactions[1].Balance + transactionResponseFirst.Transactions[0].Amount)
		assert.Equal(t, updatedAccountFirst.Balance, generatedBalanceFromPreviousTransaction)
		assert.Len(t, transactionResponseFirst.Transactions, 7)

		// Updating last transaction
		updateFields := GenerateTransactionFields(account.ID, person.ID)

		updatedTransaction, err := UpdateLastTransaction(transactionResponseFirst.Transactions[0].ID, updateFields)
		assert.Nil(t, err)

		updatedAccountSecond, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)

		transactionResponseSecond, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)

		assert.Equal(t, transactionResponseSecond.Transactions[0], updatedTransaction)
		assert.Equal(t, updatedAccountSecond.Balance, transactionResponseSecond.Transactions[0].Balance)
		generatedBalanceFromPreviousTransactionSecond := utility.RoundToTwoDecimalPlaces(transactionResponseSecond.Transactions[1].Balance + transactionResponseSecond.Transactions[0].Amount)
		assert.Equal(t, updatedAccountSecond.Balance, generatedBalanceFromPreviousTransactionSecond)

		// previous balance and new balance should be different
		assert.NotEqual(t, updatedAccountFirst.Balance, updatedAccountSecond.Balance)

		// should have still 7 transactions
		assert.Len(t, transactionResponseSecond.Transactions, 7)

		// update time should be greater than creation time
		assert.Greater(t, updatedTransaction.UpdatedAt, updatedTransaction.CreatedAt)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when trying to update a transaction that is not the last one", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(7)
		for i, v := range amounts {
			personId := person.ID
			if i%3 == 0 {
				personId = uuid.UUID{}
			}
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}

		transactionResponse, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)

		updateFields := GenerateTransactionFields(account.ID, person.ID)

		_, err = UpdateLastTransaction(transactionResponse.Transactions[1].ID, updateFields)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR003, err.Error())
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when trying to update unexisting transaction", func(t *testing.T) {
		updateFields := GenerateTransactionFields(account.ID, person.ID)

		_, err := UpdateLastTransaction(uuid.UUID{}, updateFields)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Error when updating transaction that generates negative balance", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(1)
		for i, v := range amounts {
			personId := person.ID
			if i%3 == 0 {
				personId = uuid.UUID{}
			}
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}

		transactionResponse, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)

		updateFields := GenerateTransactionFields(account.ID, person.ID)
		updateFields.Amount = -1 * (transactionResponse.Transactions[0].Balance + 1)

		_, err = UpdateLastTransaction(transactionResponse.Transactions[0].ID, updateFields)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR002, err.Error())

	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when deleting last transaction with no transactions", func(t *testing.T) {
		_, err := DeleteLastTransaction()
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	// at the end of all transactions services tests
	money_accounts.DeleteAllMoneyAccounts()
	persons.DeleteAllPersons()
}
