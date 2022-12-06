package money_accounts

import (
	"path/filepath"
	"testing"

	"github.com/grabielcruz/transportation_back/database"
)

func TestMoneyAccountsHandlers(t *testing.T) {
	envPath := filepath.Clean("../.env_test")
	sqlPath := filepath.Clean("../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()

	// r := routes.SetupRoutes()

	t.Run("Get empty slice of accounts initially", func(t *testing.T) {
		// w := httptest.NewRecorder()
		// req, _ := http.NewRequest(http.MethodGet, "/money_accounts", nil)
		// r.ServeHTTP(w, req)

		// assert.Equal(t, http.StatusOK, w.Code)
		// assert.Equal(t, "pong", w.Body.String())
	})
}
