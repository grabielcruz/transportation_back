package transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/config"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
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

	account, err := money_accounts.CreateMoneyAccount(money_accounts.GenerateAccountFields())
	assert.Nil(t, err)
	person, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)

	t.Run("Get a transaction response with zero transactions initially", func(t *testing.T) {
		url := fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		transactionsResponse := TransationResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &transactionsResponse)
		assert.Nil(t, err)
		assert.Len(t, transactionsResponse.Transactions, 0)
	})

	t.Run("Create a transaction without a person", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID, uuid.UUID{})
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transactions", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newTransaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &newTransaction)
		assert.Nil(t, err)
		assert.Equal(t, fields.AccountId, newTransaction.AccountId)
		assert.Equal(t, fields.Amount, newTransaction.Amount)
		assert.Equal(t, fields.Description, newTransaction.Description)

		transations, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Len(t, transations.Transactions, 1)
		assert.Equal(t, 1, transations.Count)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction with a person", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID, person.ID)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transactions", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newTransaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &newTransaction)
		assert.Nil(t, err)
		assert.Equal(t, fields.AccountId, newTransaction.AccountId)
		assert.Equal(t, fields.Amount, newTransaction.Amount)
		assert.Equal(t, fields.Description, newTransaction.Description)

		transations, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Len(t, transations.Transactions, 1)
		assert.Equal(t, 1, transations.Count)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when creating a transaction with an unexisting account", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(uuid.UUID{}, uuid.UUID{})
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transactions", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.TR001, errResponse.Error)
		assert.Equal(t, "TR001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when create a transaction with invalid json fields", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := generateBadTransactionFields(account.ID, person.ID)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transactions", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.UM001, errResponse.Error)
		assert.Equal(t, "UM001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when create a transaction with bad ids", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := generateBadTransactionFieldsWithBadIds(account.ID, person.ID)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transactions", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.UM001, errResponse.Error)
		assert.Equal(t, "UM001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when generating negative balance", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID, person.ID)
		fields.Amount *= -1
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transactions", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.TR002, errResponse.Error)
		assert.Equal(t, "TR002", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when sending empty description", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID, person.ID)
		fields.Description = ""
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transactions", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, "Transaction should have a description", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when sending zero balance", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID, person.ID)
		fields.Amount = float64(0)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transactions", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, "Amount should be greater than zero", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction and get it in paginated response", func(t *testing.T) {
		fields := GenerateTransactionFields(account.ID, person.ID)
		newTransaction, err := CreateTransaction(fields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var transactionsResponse TransationResponse
		err = json.Unmarshal(w.Body.Bytes(), &transactionsResponse)
		assert.Nil(t, err)
		assert.Len(t, transactionsResponse.Transactions, 1)
		assert.Equal(t, config.Offset, transactionsResponse.Offset)
		assert.Equal(t, config.Limit, transactionsResponse.Limit)
		assert.Equal(t, 1, transactionsResponse.Count)
		assert.Equal(t, newTransaction.ID, transactionsResponse.Transactions[0].ID)
		assert.Equal(t, newTransaction.Balance, transactionsResponse.Transactions[0].Balance)
		assert.Equal(t, newTransaction.Amount, transactionsResponse.Transactions[0].Amount)
		assert.Equal(t, newTransaction.AccountId, transactionsResponse.Transactions[0].AccountId)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create a transaction and get it with single response", func(t *testing.T) {
		fields := GenerateTransactionFields(account.ID, person.ID)
		newTransaction, err := CreateTransaction(fields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transaction/%v", newTransaction.ID.String())
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		transaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &transaction)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.ID, transaction.ID)
		assert.Equal(t, newTransaction.AccountId, transaction.AccountId)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, newTransaction.Balance, transaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt.UTC(), transaction.CreatedAt.UTC())
		assert.Equal(t, newTransaction.UpdatedAt.UTC(), transaction.UpdatedAt.UTC())
		assert.Equal(t, newTransaction.Date.UTC(), transaction.Date.UTC())
		assert.Equal(t, newTransaction.Description, transaction.Description)
		assert.Equal(t, newTransaction.PersonId, transaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, transaction.PersonName)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when getting unexisting transaction", func(t *testing.T) {
		url := fmt.Sprintf("/transaction/%v", (uuid.UUID{}).String())
		req, err := http.NewRequest(http.MethodGet, url, nil)
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

	t.Run("Generate 21 transactions and get them paginated", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(21)
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

		url := fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		transactionsResponse := TransationResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &transactionsResponse)
		assert.Nil(t, err)
		// the first transaction should be the last executed one
		assert.Equal(t, updatedAccount.Balance, transactionsResponse.Transactions[0].Balance)
		assert.Len(t, transactionsResponse.Transactions, config.Limit)
		assert.Equal(t, transactionsResponse.Count, 21)
		assert.Equal(t, transactionsResponse.Limit, config.Limit)
		assert.Equal(t, transactionsResponse.Offset, config.Offset)

		// offset = 10
		url = fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), config.Limit, 10)
		req, err = http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		transactionsResponse = TransationResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &transactionsResponse)
		assert.Nil(t, err)

		assert.Len(t, transactionsResponse.Transactions, config.Limit)
		assert.Equal(t, transactionsResponse.Count, 21)
		assert.Equal(t, transactionsResponse.Limit, config.Limit)
		assert.Equal(t, transactionsResponse.Offset, 10)

		// offset = 20
		url = fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), config.Limit, 20)
		req, err = http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		transactionsResponse = TransationResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &transactionsResponse)
		assert.Nil(t, err)

		assert.Len(t, transactionsResponse.Transactions, 1)
		assert.Equal(t, transactionsResponse.Count, 21)
		assert.Equal(t, transactionsResponse.Limit, config.Limit)
		assert.Equal(t, transactionsResponse.Offset, 20)
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

		buf := bytes.Buffer{}
		updateFields := GenerateTransactionFields(account.ID, person.ID)
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v", transactionResponseFirst.Transactions[0].ID)
		req, err := http.NewRequest(http.MethodPatch, url, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		updatedTransaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &updatedTransaction)
		assert.Nil(t, err)

		updatedAccountSecond, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)

		transactionResponseSecond, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)

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

		buf := bytes.Buffer{}
		updateFields := GenerateTransactionFields(account.ID, person.ID)
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v", transactionResponse.Transactions[1].ID)
		req, err := http.NewRequest(http.MethodPatch, url, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.TR003, errResponse.Error)
		assert.Equal(t, "TR003", errResponse.Code)

	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when trying to update unexisting transaction", func(t *testing.T) {
		buf := bytes.Buffer{}
		updateFields := GenerateTransactionFields(account.ID, person.ID)
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v", uuid.UUID{})
		req, err := http.NewRequest(http.MethodPatch, url, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.TR004, errResponse.Error)
		assert.Equal(t, "TR004", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when sending bad id", func(t *testing.T) {
		buf := bytes.Buffer{}
		updateFields := GenerateTransactionFields(account.ID, person.ID)
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v", "abcde")
		req, err := http.NewRequest(http.MethodPatch, url, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 5", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when trying to update a transaction that generates a negative balance", func(t *testing.T) {
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

		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)

		buf := bytes.Buffer{}
		updateFields := GenerateTransactionFields(account.ID, person.ID)
		updateFields.Amount = -1 * (updatedAccount.Balance + 1)
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v", transactionResponse.Transactions[0].ID)
		req, err := http.NewRequest(http.MethodPatch, url, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.TR002, errResponse.Error)
		assert.Equal(t, "TR002", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when updating transaction with empty description", func(t *testing.T) {
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

		buf := bytes.Buffer{}
		updateFields := GenerateTransactionFields(account.ID, person.ID)
		updateFields.Description = ""
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v", transactionResponse.Transactions[0].ID)
		req, err := http.NewRequest(http.MethodPatch, url, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Transaction should have a description", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when updating last transaction with zero amount", func(t *testing.T) {
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

		buf := bytes.Buffer{}
		updateFields := GenerateTransactionFields(account.ID, person.ID)
		updateFields.Amount = float64(0)
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v", transactionResponse.Transactions[0].ID)
		req, err := http.NewRequest(http.MethodPatch, url, &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Amount should be greater than zero", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when deleting last transaction with no transactions", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/transactions", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.TR004, errResponse.Error)
		assert.Equal(t, "TR004", errResponse.Code)
	})

	// at the end of all transactions services tests
	money_accounts.DeleteAllMoneyAccounts()
	persons.DeleteAllPersons()
}
