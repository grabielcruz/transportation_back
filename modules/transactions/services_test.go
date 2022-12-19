package transactions

import (
	"path/filepath"
	"testing"

	"github.com/grabielcruz/transportation_back/database"
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

	})
}
