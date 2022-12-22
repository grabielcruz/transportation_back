package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/grabielcruz/transportation_back/database"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestTransactionsHandlers(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	router := httprouter.New()
	Routes(router)
	w := httptest.NewRecorder()

	account := money_accounts.CreateMoneyAccount(money_accounts.GenerateAccountFields())
	// person := persons.CreatePerson(persons.GeneratePersonFields())

	t.Run("Get a transaction response with zero transactions initially", func(t *testing.T) {
		url := fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), Limit, Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		transactionsResponse := TransationResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &transactionsResponse)
		assert.Nil(t, err)
		assert.Len(t, transactionsResponse.Transactions, 0)
	})
}
