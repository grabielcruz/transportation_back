package currencies

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
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

	t.Run("Get slice of two currencies initially", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/currencies", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		currencies := []string{}
		err = json.Unmarshal(w.Body.Bytes(), &currencies)
		assert.Nil(t, err)
		assert.Len(t, currencies, 2)
	})

	t.Run("Can create a currency", func(t *testing.T) {
		buf := bytes.Buffer{}
		newCurrency := "ABC"
		req, err := http.NewRequest(http.MethodPost, "/currencies/"+newCurrency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		createdCurrency := ""
		err = json.Unmarshal(w.Body.Bytes(), &createdCurrency)
		assert.Nil(t, err)
		assert.Equal(t, newCurrency, createdCurrency)

		// get currencies
		req2, err := http.NewRequest(http.MethodGet, "/currencies", nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		currencies := []string{}
		err = json.Unmarshal(w2.Body.Bytes(), &currencies)
		assert.Nil(t, err)
		assert.Len(t, currencies, 3)
	})

	resetCurrencies()

	t.Run("Error when creating currency with bad format", func(t *testing.T) {
		buf := bytes.Buffer{}
		newCurrency := "bad_currency"
		req, err := http.NewRequest(http.MethodPost, "/currencies/"+newCurrency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.CU002, errResponse.Error)
		assert.Equal(t, "CU002", errResponse.Code)
	})

	t.Run("Error when creating repeated currency", func(t *testing.T) {
		buf := bytes.Buffer{}
		newCurrency := "VED"
		req, err := http.NewRequest(http.MethodPost, "/currencies/"+newCurrency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.CU003, errResponse.Error)
		assert.Equal(t, "CU003", errResponse.Code)
	})

	t.Run("Error when creating bad currency", func(t *testing.T) {
		buf := bytes.Buffer{}
		newCurrency := ""
		req, err := http.NewRequest(http.MethodPost, "/currencies/"+newCurrency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Can create a currency, then delete it", func(t *testing.T) {
		buf := bytes.Buffer{}
		newCurrency := "ABC"
		req, err := http.NewRequest(http.MethodPost, "/currencies/"+newCurrency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		createdCurrency := ""
		err = json.Unmarshal(w.Body.Bytes(), &createdCurrency)
		assert.Nil(t, err)
		assert.Equal(t, newCurrency, createdCurrency)

		req2, err := http.NewRequest(http.MethodDelete, "/currencies/"+createdCurrency, nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		deletedCurrency := ""
		err = json.Unmarshal(w2.Body.Bytes(), &deletedCurrency)
		assert.Nil(t, err)

		assert.Equal(t, createdCurrency, deletedCurrency)

		// checking
		currencies := GetCurrencies()
		assert.Len(t, currencies, 2)
		assert.Equal(t, "VED", currencies[0])
		assert.Equal(t, "USD", currencies[1])
	})

	resetCurrencies()

	t.Run("Error when deleting unexisting currency", func(t *testing.T) {
		buf := bytes.Buffer{}
		currency := "WWW"
		req, err := http.NewRequest(http.MethodDelete, "/currencies/"+currency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("Error when trying to delete VED currency", func(t *testing.T) {
		buf := bytes.Buffer{}
		currency := "VED"
		req, err := http.NewRequest(http.MethodDelete, "/currencies/"+currency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.CU001, errResponse.Error)
		assert.Equal(t, "CU001", errResponse.Code)
	})

	t.Run("Error when trying to delete USD currency", func(t *testing.T) {
		buf := bytes.Buffer{}
		currency := "USD"
		req, err := http.NewRequest(http.MethodDelete, "/currencies/"+currency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.CU001, errResponse.Error)
		assert.Equal(t, "CU001", errResponse.Code)
	})

	t.Run("Error when trying to delete currency associated with a money account", func(t *testing.T) {
		createdCurrency, err := CreateCurrency("ABC")
		assert.Nil(t, err)

		accountsFields := money_accounts.GenerateAccountFields()
		accountsFields.Currency = createdCurrency
		newMoneyAccount, err := money_accounts.CreateMoneyAccount(accountsFields)
		assert.Nil(t, err)
		assert.Equal(t, newMoneyAccount.Currency, createdCurrency)

		buf := bytes.Buffer{}
		req, err := http.NewRequest(http.MethodDelete, "/currencies/"+createdCurrency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.CU004, errResponse.Error)
		assert.Equal(t, "CU004", errResponse.Code)
	})

	t.Run("Error when deleting zero currency", func(t *testing.T) {
		buf := bytes.Buffer{}
		currency := "000"
		req, err := http.NewRequest(http.MethodDelete, "/currencies/"+currency, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.CU002, errResponse.Error)
		assert.Equal(t, "CU002", errResponse.Code)
	})

}
