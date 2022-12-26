package money_accounts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestMoneyAccountsHandlers(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	router := httprouter.New()
	Routes(router)
	w := httptest.NewRecorder()

	t.Run("Get empty slice of accounts initially", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/money_accounts", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		accounts := []MoneyAccount{}
		err = json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.Nil(t, err)
		assert.Len(t, accounts, 0)
	})

	t.Run("Create one money account", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateAccountFields()
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/money_accounts", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		createdAccount := MoneyAccount{}
		err = json.Unmarshal(w.Body.Bytes(), &createdAccount)
		assert.Nil(t, err)
		assert.Equal(t, fields.Name, createdAccount.Name)
		assert.Equal(t, fields.Details, createdAccount.Details)
		assert.Equal(t, fields.Currency, createdAccount.Currency)
	})

	DeleteAllMoneyAccounts()

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

	DeleteAllMoneyAccounts()

	t.Run("Error when sending invalid json when creating account", func(t *testing.T) {
		buf := bytes.Buffer{}
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
		assert.Equal(t, errors_handler.UM001, errResponse.Error)
		assert.Equal(t, "UM001", errResponse.Code)
	})

	t.Run("Error when sending bad fields on creating a person", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := MoneyAccountFields{}
		err := json.NewEncoder(&buf).Encode(fields)
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
		assert.Equal(t, "Name is required", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
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
		assert.Equal(t, fields.Details, account.Details)
		assert.Equal(t, fields.Currency, account.Currency)
	})

	DeleteAllMoneyAccounts()

	t.Run("Get error when sending bad id", func(t *testing.T) {
		badId := utility.GetRandomString(10)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts/"+badId, nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 10", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Get error when sending uregistered id", func(t *testing.T) {
		wantedId := uuid.UUID{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/money_accounts/"+wantedId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("It should create and update one money account", func(t *testing.T) {
		createFields := GenerateAccountFields()
		wantedId := CreateMoneyAccount(createFields).ID
		buf := bytes.Buffer{}
		updateFields := GenerateAccountFields()
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/money_accounts/"+wantedId.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		updatedAccount := MoneyAccount{}
		err = json.Unmarshal(w.Body.Bytes(), &updatedAccount)
		assert.Nil(t, err)
		assert.Equal(t, updateFields.Name, updatedAccount.Name)
		assert.Equal(t, updateFields.Details, updatedAccount.Details)
		assert.Equal(t, updateFields.Currency, updatedAccount.Currency)
		assert.Greater(t, updatedAccount.UpdatedAt, updatedAccount.CreatedAt)
	})

	DeleteAllMoneyAccounts()

	t.Run("Error when sending bad id", func(t *testing.T) {
		badId := utility.GetRandomString(10)
		w := httptest.NewRecorder()
		buf := bytes.Buffer{}
		updateFields := GenerateAccountFields()
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		req, err := http.NewRequest(http.MethodPatch, "/money_accounts/"+badId, &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		errResponse := errors_handler.ErrorResponse{}

		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 10", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Error when sending unregistered id when patching", func(t *testing.T) {
		wantedId := uuid.UUID{}
		w := httptest.NewRecorder()
		buf := bytes.Buffer{}
		updateFields := GenerateAccountFields()
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		req, err := http.NewRequest(http.MethodPatch, "/money_accounts/"+wantedId.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		errResponse := errors_handler.ErrorResponse{}

		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("Error when sending bad json on updating account", func(t *testing.T) {
		buf := bytes.Buffer{}
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
		assert.Equal(t, errors_handler.UM001, errResponse.Error)
		assert.Equal(t, "UM001", errResponse.Code)
	})

	t.Run("Error when sending bad fields on updating account", func(t *testing.T) {
		buf := bytes.Buffer{}
		wantedId := uuid.UUID{}
		badFields := MoneyAccountFields{}
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
		assert.Equal(t, "Name is required", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("It should create an account and delete it", func(t *testing.T) {
		fields := GenerateAccountFields()
		newId := CreateMoneyAccount(fields).ID

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/money_accounts/"+newId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deletedId := common.ID{}
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &deletedId)
		assert.Nil(t, err)
		assert.Equal(t, newId, deletedId.ID)

		deletedAccount, err := GetOneMoneyAccount(newId)
		assert.Equal(t, deletedAccount.ID, uuid.UUID{})
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	DeleteAllMoneyAccounts()

	t.Run("it should send error when sending bad id", func(t *testing.T) {
		newId := utility.GetRandomString(10)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/money_accounts/"+newId, nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, "invalid UUID length: 10", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("it should send error when trying to delete unregistered account", func(t *testing.T) {
		newId := uuid.UUID{}

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/money_accounts/"+newId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})
}
