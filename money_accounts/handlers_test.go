package money_accounts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestMoneyAccountsHandlers(t *testing.T) {
	envPath := filepath.Clean("../.env_test")
	sqlPath := filepath.Clean("../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	router := httprouter.New()
	Routes(router)

	t.Run("Get empty slice of accounts initially", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
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

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var createdAccount MoneyAccount
		err = json.Unmarshal(w.Body.Bytes(), &createdAccount)
		assert.Nil(t, err)
		assert.Equal(t, fields.Name, createdAccount.Name)
		assert.Equal(t, fields.IsCash, createdAccount.IsCash)
		assert.Equal(t, fields.Currency, createdAccount.Currency)
	})

	deleteAllMoneyAccounts()

	t.Run("Create three money accounts and get an slice of accounts", func(t *testing.T) {
		CreateMoneyAccount(GenerateAccountFields())
		CreateMoneyAccount(GenerateAccountFields())
		CreateMoneyAccount(GenerateAccountFields())
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		accounts := []MoneyAccount{}

		err = json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.Nil(t, err)
		assert.Len(t, accounts, 3)
	})

	deleteAllMoneyAccounts()

	t.Run("Error when sending dummy json when creating account", func(t *testing.T) {
		var buf bytes.Buffer
		badFields := generateBadAccountFields()
		err := json.NewEncoder(&buf).Encode(badFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/money_accounts", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, "Invalid data type", errResponse.Error)
	})

	t.Run("Create one money account and get it", func(t *testing.T) {
		fields := GenerateAccountFields()
		wantedId := CreateMoneyAccount(fields).ID
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts/"+wantedId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var account MoneyAccount

		err = json.Unmarshal(w.Body.Bytes(), &account)
		assert.Nil(t, err)
		assert.Equal(t, fields.Name, account.Name)
		assert.Equal(t, fields.IsCash, account.IsCash)
		assert.Equal(t, fields.Currency, account.Currency)
	})

	deleteAllMoneyAccounts()

	t.Run("Get error when sending wrong id", func(t *testing.T) {
		wantedId := "abcdefg"
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts/"+wantedId, nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var errResponse errors_handler.ErrorResponse

		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 7", errResponse.Error)
	})

	t.Run("Get error when sending unexisting id", func(t *testing.T) {
		wantedId := uuid.UUID{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts/"+wantedId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var errResponse errors_handler.ErrorResponse

		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "sql: no rows in result set", errResponse.Error)
	})

	t.Run("It Should create and update one money account", func(t *testing.T) {
		createFields := GenerateAccountFields()
		wantedId := CreateMoneyAccount(createFields).ID
		var buf bytes.Buffer
		updateFields := GenerateAccountFields()
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/money_accounts/"+wantedId.String(), &buf)
		assert.Nil(t, err)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var updatedAccount MoneyAccount
		err = json.Unmarshal(w.Body.Bytes(), &updatedAccount)
		assert.Nil(t, err)

		assert.Equal(t, updateFields.Name, updatedAccount.Name)
		assert.Equal(t, updateFields.IsCash, updatedAccount.IsCash)
		assert.Equal(t, updateFields.Currency, updatedAccount.Currency)
	})

	deleteAllMoneyAccounts()

	t.Run("Error when sending wrong id", func(t *testing.T) {
		wantedId := "abcdefg"
		w := httptest.NewRecorder()
		var buf bytes.Buffer
		updateFields := GenerateAccountFields()
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		req, err := http.NewRequest(http.MethodPatch, "/money_accounts/"+wantedId, &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var errResponse errors_handler.ErrorResponse

		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 7", errResponse.Error)
	})

	t.Run("Error when sending unexisting id when patching", func(t *testing.T) {
		wantedId := uuid.UUID{}
		w := httptest.NewRecorder()
		var buf bytes.Buffer
		updateFields := GenerateAccountFields()
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		req, err := http.NewRequest(http.MethodPatch, "/money_accounts/"+wantedId.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var errResponse errors_handler.ErrorResponse

		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "sql: no rows in result set", errResponse.Error)
	})

	t.Run("Error when sending dummy json when updating account", func(t *testing.T) {
		var buf bytes.Buffer
		wantedId := uuid.UUID{}
		badFields := generateBadAccountFields()
		err := json.NewEncoder(&buf).Encode(badFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/money_accounts/"+wantedId.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, "Invalid data type", errResponse.Error)
	})

	//todo error when sending patch without body
}
