package money_accounts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/stretchr/testify/assert"
)

func TestMoneyAccountsHandlers(t *testing.T) {
	envPath := filepath.Clean("../.env_test")
	sqlPath := filepath.Clean("../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	r := gin.Default()
	Routes(r)

	t.Run("Get empty slice of accounts initially", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts", nil)
		assert.Nil(t, err)

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var accounts []MoneyAccount

		err = json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.Nil(t, err)
		assert.Len(t, accounts, 0)
	})

	t.Run("Create one money account", func(t *testing.T) {
		var buf bytes.Buffer
		fields := GenerateAccountFields()
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/money_accounts", &buf)
		assert.Nil(t, err)

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var createdAccount MoneyAccount
		err = json.Unmarshal(w.Body.Bytes(), &createdAccount)
		assert.Nil(t, err)
		assert.Equal(t, createdAccount.Name, fields.Name)
		assert.Equal(t, createdAccount.IsCash, fields.IsCash)
		assert.Equal(t, createdAccount.Currency, fields.Currency)
	})

	deleteAllMoneyAccounts()

	t.Run("Create three money accounts and get an slice of accounts", func(t *testing.T) {
		CreateMoneyAccount(GenerateAccountFields())
		CreateMoneyAccount(GenerateAccountFields())
		CreateMoneyAccount(GenerateAccountFields())
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts", nil)
		assert.Nil(t, err)

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		accounts := []MoneyAccount{}

		err = json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.Nil(t, err)
		assert.Len(t, accounts, 3)
	})

	deleteAllMoneyAccounts()

	t.Run("Error when sending dummy json", func(t *testing.T) {
		var buf bytes.Buffer
		dummy := utility.GenerateDummyData()
		err := json.NewEncoder(&buf).Encode(dummy)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/money_accounts", &buf)
		assert.Nil(t, err)

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
	})
}
