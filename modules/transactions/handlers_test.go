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

	account := money_accounts.CreateMoneyAccount(money_accounts.GenerateAccountFields())
	person := persons.CreatePerson(persons.GeneratePersonFields())

	t.Run("Get a transaction response with zero transactions initially", func(t *testing.T) {
		url := fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), Limit, Offset)
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

		transations, err := GetTransactions(account.ID, Limit, Offset)
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

		transations, err := GetTransactions(account.ID, Limit, Offset)
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

	t.Run("Create one transaction and get it", func(t *testing.T) {
		fields := GenerateTransactionFields(account.ID, person.ID)
		newTransaction, err := CreateTransaction(fields)
		assert.Nil(t, err)

		url := fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), Limit, Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var transactionsResponse TransationResponse
		err = json.Unmarshal(w.Body.Bytes(), &transactionsResponse)
		assert.Nil(t, err)
		assert.Len(t, transactionsResponse.Transactions, 1)
		assert.Equal(t, Offset, transactionsResponse.Offset)
		assert.Equal(t, Limit, transactionsResponse.Limit)
		assert.Equal(t, 1, transactionsResponse.Count)
		assert.Equal(t, newTransaction.ID, transactionsResponse.Transactions[0].ID)
		assert.Equal(t, newTransaction.Balance, transactionsResponse.Transactions[0].Balance)
		assert.Equal(t, newTransaction.Amount, transactionsResponse.Transactions[0].Amount)
		assert.Equal(t, newTransaction.AccountId, transactionsResponse.Transactions[0].AccountId)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

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

		url := fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), Limit, Offset)
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
		assert.Len(t, transactionsResponse.Transactions, Limit)
		assert.Equal(t, transactionsResponse.Count, 21)
		assert.Equal(t, transactionsResponse.Limit, Limit)
		assert.Equal(t, transactionsResponse.Offset, Offset)

		// offset = 10
		url = fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), Limit, 10)
		req, err = http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		transactionsResponse = TransationResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &transactionsResponse)
		assert.Nil(t, err)

		assert.Len(t, transactionsResponse.Transactions, Limit)
		assert.Equal(t, transactionsResponse.Count, 21)
		assert.Equal(t, transactionsResponse.Limit, Limit)
		assert.Equal(t, transactionsResponse.Offset, 10)

		// offset = 20
		url = fmt.Sprintf("/transactions/%v?limit=%v&offset=%v", account.ID.String(), Limit, 20)
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
		assert.Equal(t, transactionsResponse.Limit, Limit)
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

		transactionResponseFirst, err := GetTransactions(account.ID, Limit, Offset)
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

		transactionResponseSecond, err := GetTransactions(account.ID, Limit, Offset)
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

		transactionResponse, err := GetTransactions(account.ID, Limit, Offset)
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

	t.Run("Error when trying to update unexisting transaction with empty database", func(t *testing.T) {
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

		transactionResponse, err := GetTransactions(account.ID, Limit, Offset)
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

		transactionResponse, err := GetTransactions(account.ID, Limit, Offset)
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

		transactionResponse, err := GetTransactions(account.ID, Limit, Offset)
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
	resetTransactions()

	t.Run("Should create two transactions, delete the last one and get it from trashed transactions", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(2)
		for _, v := range amounts {
			personId := person.ID
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}

		// delete transaction
		req, err := http.NewRequest(http.MethodDelete, "/transactions", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deletedTransaction := TrashedTransaction{}
		err = json.Unmarshal(w.Body.Bytes(), &deletedTransaction)
		assert.Nil(t, err)
		/////////////////

		transactions, err := GetTransactions(account.ID, Limit, Offset)
		assert.Nil(t, err)

		assert.Len(t, transactions.Transactions, 1)

		// get trashed transactions
		req2, err := http.NewRequest(http.MethodGet, "/trashed_transactions", nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w.Code)

		trashedTransactions := []TrashedTransaction{}
		err = json.Unmarshal(w2.Body.Bytes(), &trashedTransactions)
		assert.Nil(t, err)
		/////////////////

		assert.Equal(t, deletedTransaction.ID, trashedTransactions[0].ID)
		assert.Equal(t, deletedTransaction.AccountId, trashedTransactions[0].AccountId)
		assert.Equal(t, deletedTransaction.PersonId, trashedTransactions[0].PersonId)
		assert.Equal(t, deletedTransaction.Amount, trashedTransactions[0].Amount)
		assert.Equal(t, deletedTransaction.Date.UTC(), trashedTransactions[0].Date.UTC())
		assert.Equal(t, deletedTransaction.Description, trashedTransactions[0].Description)
		assert.Equal(t, deletedTransaction.CreatedAt.UTC(), trashedTransactions[0].CreatedAt.UTC())
		assert.Equal(t, deletedTransaction.UpdatedAt.UTC(), trashedTransactions[0].UpdatedAt.UTC())
		assert.Equal(t, deletedTransaction.DeletedAt.UTC(), trashedTransactions[0].DeletedAt.UTC())
		assert.Greater(t, trashedTransactions[0].DeletedAt.UTC(), trashedTransactions[0].CreatedAt.UTC())

		// it should update account's balance, not a problem if the last transaction is deleted
		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)
		assert.Equal(t, updatedAccount.Balance, transactions.Transactions[0].Balance)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	resetTransactions()

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

	t.Run("Should create two transactions, delete the last one and restore it", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(2)
		for _, v := range amounts {
			personId := person.ID
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}

		// delete last transaction
		req, err := http.NewRequest(http.MethodDelete, "/transactions", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deletedTransaction := TrashedTransaction{}
		err = json.Unmarshal(w.Body.Bytes(), &deletedTransaction)
		assert.Nil(t, err)
		///////////////////////////

		// restore last transaction
		buf := bytes.Buffer{}
		req2, err := http.NewRequest(http.MethodPost, "/trashed_transactions/"+deletedTransaction.ID.String(), &buf)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w.Code)

		restoredTransaction := Transaction{}
		err = json.Unmarshal(w2.Body.Bytes(), &restoredTransaction)
		assert.Nil(t, err)
		///////////////////////////////

		transactions, err := GetTransactions(account.ID, Limit, Offset)
		assert.Nil(t, err)
		assert.Len(t, transactions.Transactions, 2)

		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)
		assert.Equal(t, updatedAccount.Balance, transactions.Transactions[0].Balance)

		// get trashed transactions
		req3, err := http.NewRequest(http.MethodGet, "/trashed_transactions", nil)
		assert.Nil(t, err)

		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w.Code)

		trashedTransactions := []TrashedTransaction{}
		err = json.Unmarshal(w3.Body.Bytes(), &trashedTransactions)
		assert.Nil(t, err)
		/////////////////

		assert.Len(t, trashedTransactions, 0)

		assert.Equal(t, deletedTransaction.ID, restoredTransaction.ID)
		assert.Equal(t, deletedTransaction.AccountId, restoredTransaction.AccountId)
		assert.Equal(t, deletedTransaction.PersonId, restoredTransaction.PersonId)
		assert.Equal(t, deletedTransaction.Amount, restoredTransaction.Amount)
		assert.Equal(t, deletedTransaction.Date, restoredTransaction.Date)
		assert.Equal(t, deletedTransaction.Description, restoredTransaction.Description)
		assert.Less(t, deletedTransaction.DeletedAt, restoredTransaction.CreatedAt)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	resetTransactions()

	t.Run("Error when restoring unexisting transaction", func(t *testing.T) {
		buf := bytes.Buffer{}
		req, err := http.NewRequest(http.MethodPost, "/trashed_transactions/"+uuid.UUID{}.String(), &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.TR011, errResponse.Error)
		assert.Equal(t, "TR011", errResponse.Code)
	})

	t.Run("Error when restoring transaction with bad id", func(t *testing.T) {
		buf := bytes.Buffer{}
		req, err := http.NewRequest(http.MethodPost, "/trashed_transactions/"+"bad_id", &buf)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 6", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Error when restoring transaction that generates negative balance", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID, person.ID)
		transactionFields.Amount = float64(100)
		_, err := CreateTransaction(transactionFields)
		assert.Nil(t, err)

		transactionFields2 := GenerateTransactionFields(account.ID, person.ID)
		transactionFields2.Amount = float64(-100)
		_, err = CreateTransaction(transactionFields2)
		assert.Nil(t, err)

		// delete last transaction
		req, err := http.NewRequest(http.MethodDelete, "/transactions", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deletedTransaction := TrashedTransaction{}
		err = json.Unmarshal(w.Body.Bytes(), &deletedTransaction)
		assert.Nil(t, err)
		///////////////////////////

		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)

		transactionFields3 := GenerateTransactionFields(account.ID, uuid.UUID{})
		transactionFields3.Amount = updatedAccount.Balance * -1 // generates zero balance on account

		_, err = CreateTransaction(transactionFields3)
		assert.Nil(t, err)

		// restore last transaction
		buf := bytes.Buffer{}
		req2, err := http.NewRequest(http.MethodPost, "/trashed_transactions/"+deletedTransaction.ID.String(), &buf)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w2.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.TR002, errResponse.Error)
		assert.Equal(t, "TR002", errResponse.Code)
		///////////////////////////////
	})

	money_accounts.ResetAccountsBalance(account.ID)
	resetTransactions()

	t.Run("Should create two transactions, delete the last one and then delete it permanently", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(2)
		for _, v := range amounts {
			personId := person.ID
			transactionFields := GenerateTransactionFields(account.ID, personId)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields)
			assert.Nil(t, err)
		}

		// delete last transaction
		req, err := http.NewRequest(http.MethodDelete, "/transactions", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		lastTransaction := TrashedTransaction{}
		err = json.Unmarshal(w.Body.Bytes(), &lastTransaction)
		assert.Nil(t, err)
		///////////////////////////

		// delete trashed transaction permanently
		req2, err := http.NewRequest(http.MethodDelete, "/trashed_transactions/"+lastTransaction.ID.String(), nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w.Code)

		permanently_deleted := TrashedTransaction{}
		err = json.Unmarshal(w.Body.Bytes(), &permanently_deleted)
		assert.Nil(t, err)
		///////////////////////////

		transactions, err := GetTransactions(account.ID, Limit, Offset)
		assert.Nil(t, err)
		assert.Len(t, transactions.Transactions, 1)

		// get trashed transactions
		req3, err := http.NewRequest(http.MethodGet, "/trashed_transactions", nil)
		assert.Nil(t, err)

		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w.Code)

		trashedTransactions := []TrashedTransaction{}
		err = json.Unmarshal(w3.Body.Bytes(), &trashedTransactions)
		assert.Nil(t, err)
		/////////////////

		assert.Len(t, trashedTransactions, 0)
		assert.Equal(t, lastTransaction.ID, permanently_deleted.ID)
		assert.Equal(t, lastTransaction.AccountId, permanently_deleted.AccountId)
		assert.Equal(t, lastTransaction.PersonId, permanently_deleted.PersonId)
		assert.Equal(t, lastTransaction.Amount, permanently_deleted.Amount)
		assert.Equal(t, lastTransaction.Date, permanently_deleted.Date)
		assert.Equal(t, lastTransaction.Description, permanently_deleted.Description)
		assert.Equal(t, lastTransaction.DeletedAt, permanently_deleted.DeletedAt)

	})

	money_accounts.ResetAccountsBalance(account.ID)
	resetTransactions()

	t.Run("Error when deleting unexisting trashed transaction", func(t *testing.T) {
		// delete trashed transaction permanently
		req, err := http.NewRequest(http.MethodDelete, "/trashed_transactions/"+uuid.UUID{}.String(), nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		///////////////////////////
		assert.Equal(t, errors_handler.TR011, errResponse.Error)
		assert.Equal(t, "TR011", errResponse.Code)
	})

	t.Run("Error when deleting trashed transaction with bad id", func(t *testing.T) {
		// delete trashed transaction permanently
		req, err := http.NewRequest(http.MethodDelete, "/trashed_transactions/"+"bad_id", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		///////////////////////////
		assert.Equal(t, "invalid UUID length: 6", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})
	// at the end of all transactions services tests
	money_accounts.DeleteAllMoneyAccounts()
	persons.DeleteAllPersons()
}
