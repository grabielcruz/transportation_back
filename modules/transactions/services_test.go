package transactions

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/stretchr/testify/assert"
)

func TestTransactionServices(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()

	t.Run("Get transaction response with zero transactions", func(t *testing.T) {
		transactions := GetTransactions(Offset, Limit)
		assert.Len(t, transactions.Transactions, 0)
		assert.Equal(t, transactions.Pagination.Count, 0)
		assert.Equal(t, transactions.Pagination.Offset, 0)
		assert.Equal(t, transactions.Pagination.Limit, 0)
	})

	t.Run("Create one transaction without person", func(t *testing.T) {
		accountFields := money_accounts.GenerateAccountFields()
		newAccount := money_accounts.CreateMoneyAccount(accountFields)
		transactionFields := GenerateTransactionFields(newAccount.ID, uuid.UUID{})
		newTransaction, err := CreateTransaction(transactionFields)
		assert.Nil(t, err)
		updatedAccount, err := money_accounts.GetOneMoneyAccount(newTransaction.AccountId)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.Balance, updatedAccount.Balance)
	})
}
