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
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/bills"
	"github.com/grabielcruz/transportation_back/modules/config"
	"github.com/grabielcruz/transportation_back/modules/currencies"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/grabielcruz/transportation_back/modules/person_accounts"
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
	bills.Routes(router)
	person_accounts.Routes(router)

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

	t.Run("Create one transaction with a person", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
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

	t.Run("Create one transaction with a person and a person account, then delete the person account", func(t *testing.T) {
		// person account
		personAccountFields := person_accounts.GeneratePersonAccountFields()
		// force same currency
		personAccountFields.Currency = account.Currency
		newPersonAccount, err := person_accounts.CreatePersonAccount(person.ID, personAccountFields)
		assert.Nil(t, err)
		//

		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		fields.PersonAccountId = newPersonAccount.ID
		err = json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newTransaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &newTransaction)
		assert.Nil(t, err)

		updatedAccount, err := money_accounts.GetOneMoneyAccount(newTransaction.AccountId)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.Balance, updatedAccount.Balance)
		assert.Equal(t, newTransaction.AccountId, updatedAccount.ID)
		assert.Equal(t, newTransaction.PersonName, person.Name)
		assert.Equal(t, newTransaction.PersonId, person.ID)
		// persons account
		assert.Equal(t, newTransaction.PersonAccountId, newPersonAccount.ID)
		assert.Equal(t, newTransaction.PersonAccountName, newPersonAccount.Name)
		assert.Equal(t, newTransaction.PersonAccountDescription, newPersonAccount.Description)
		assert.Equal(t, newTransaction.Currency, newPersonAccount.Currency)

		// delete person account
		w2 := httptest.NewRecorder()
		req2, err := http.NewRequest(http.MethodDelete, "/one_person_account/"+newPersonAccount.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		deletedId := common.ID{}
		err = json.Unmarshal(w2.Body.Bytes(), &deletedId)
		assert.Nil(t, err)
		assert.Equal(t, newPersonAccount.ID, deletedId.ID)
		obtainedTransaction, err := GetTransaction(newTransaction.ID)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.ID, obtainedTransaction.ID)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	person_accounts.DeleteAllPersonAccounts()
	deleteAllTransactions()

	// ESTE ES PERSON_ACOUNT
	t.Run("Error when creating a transaction with an unexisting person account different than zero", func(t *testing.T) {
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.PersonAccountId = randId

		buf := bytes.Buffer{}
		err = json.NewEncoder(&buf).Encode(transactionFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.PA002, errResponse.Error)
		assert.Equal(t, "PA002", errResponse.Code)
	})

	// ESTE ES PERSON_ACOUNT
	t.Run("Error when creating a transaction with an person account with currency mismatch", func(t *testing.T) {
		// person account
		personAccountFields := person_accounts.GeneratePersonAccountFields()
		// force different currency
		newCurrency, err := currencies.CreateCurrency("ABC")
		assert.Nil(t, err)
		personAccountFields.Currency = newCurrency
		newPersonAccount, err := person_accounts.CreatePersonAccount(person.ID, personAccountFields)
		assert.Nil(t, err)
		//
		fields := GenerateTransactionFields(account.ID)
		fields.PersonAccountId = newPersonAccount.ID

		buf := bytes.Buffer{}
		err = json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.TR011, errResponse.Error)
		assert.Equal(t, "TR011", errResponse.Code)
	})

	t.Run("Error when creating a transaction with an unexisting account", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(uuid.UUID{})
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
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
		fields := generateBadTransactionFields(account.ID)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
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
		fields := generateBadTransactionFieldsWithBadIds(account.ID)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+"5555a", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, "invalid UUID length: 5", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when generating negative balance", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		fields.Amount *= -1
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
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
		fields := GenerateTransactionFields(account.ID)
		fields.Description = ""
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
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

	t.Run("Error when sending zero amount", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		fields.Amount = float64(0)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
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

	t.Run("Error when sending negative fee", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		fields.Fee = -0.1
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.TR009, errResponse.Error)
		assert.Equal(t, "TR009", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when sending fee greater than one", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		fields.Fee = 1.05
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse errors_handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.TR009, errResponse.Error)
		assert.Equal(t, "TR009", errResponse.Code)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction and get it in paginated response", func(t *testing.T) {
		fields := GenerateTransactionFields(account.ID)
		newTransaction, err := CreateTransaction(fields, person.ID, true)
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

	t.Run("Create a transaction without fee and get it with single response", func(t *testing.T) {
		// creating
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		fields.Fee = 0
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newTransaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &newTransaction)
		assert.Nil(t, err)

		// getting
		url := fmt.Sprintf("/transaction/%v", newTransaction.ID.String())
		req2, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

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

	t.Run("Create a transaction with fee and get it with single response", func(t *testing.T) {
		// creating
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		fields.Fee = 0
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newTransaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &newTransaction)
		assert.Nil(t, err)

		// getting
		url := fmt.Sprintf("/transaction/%v", newTransaction.ID.String())
		req2, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		transaction := Transaction{}
		err = json.Unmarshal(w2.Body.Bytes(), &transaction)
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

	t.Run("Error when creating transaction without a person on pending bill url", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GenerateTransactionFields(account.ID)
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transaction_to_pending_bill/"+(uuid.UUID{}).String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.TR007, errResponse.Error)
		assert.Equal(t, "TR007", errResponse.Code)
	})

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
		for _, v := range amounts {
			personId := person.ID

			transactionFields := GenerateTransactionFields(account.ID)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields, personId, true)
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

	t.Run("Error when creating a transaction without fee, delete it and then getting it", func(t *testing.T) {
		fields := GenerateTransactionFields(account.ID)
		fields.Fee = 0
		newTransaction, err := CreateTransaction(fields, person.ID, true)
		assert.Nil(t, err)

		// pending bill
		newPendingBill, err := bills.GetOneBill(newTransaction.PendingBillId)
		assert.Nil(t, err)
		assert.Equal(t, newPendingBill.ParentTransactionId, newTransaction.ID)
		assert.Equal(t, newPendingBill.Amount, newTransaction.AmountWithFee)
		assert.Equal(t, newPendingBill.Amount, newTransaction.Amount) // fee 0
		assert.Equal(t, newPendingBill.Date, newTransaction.Date)
		assert.Equal(t, newPendingBill.Description, newTransaction.Description)

		// delete
		req, err := http.NewRequest(http.MethodDelete, "/transactions", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deletedTransaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &deletedTransaction)
		assert.Nil(t, err)

		assert.Equal(t, newTransaction.ID, deletedTransaction.ID)
		assert.Equal(t, newTransaction.AccountId, deletedTransaction.AccountId)
		assert.Equal(t, newTransaction.Amount, deletedTransaction.Amount)
		assert.Equal(t, newTransaction.Balance, deletedTransaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt.UTC(), deletedTransaction.CreatedAt.UTC())
		assert.Equal(t, newTransaction.UpdatedAt.UTC(), deletedTransaction.UpdatedAt.UTC())
		assert.Equal(t, newTransaction.Date.UTC(), deletedTransaction.Date.UTC())
		assert.Equal(t, newTransaction.Description, deletedTransaction.Description)
		assert.Equal(t, newTransaction.PersonId, deletedTransaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, deletedTransaction.PersonName)

		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		// assert.Equal(t, deletedLastTransaction.PendingBillId, uuid.UUID{})
		assert.Equal(t, deletedTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, deletedTransaction.RevertBillId, uuid.UUID{})

		url := fmt.Sprintf("/transaction/%v", newTransaction.ID.String())
		req2, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w2.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)

		// pending bill also deleted
		_, err = bills.GetOneBill(newTransaction.PendingBillId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when creating a transaction with fee, delete it and then getting it", func(t *testing.T) {
		fields := GenerateTransactionFields(account.ID)
		newTransaction, err := CreateTransaction(fields, person.ID, true)
		assert.Nil(t, err)

		// pending bill
		newPendingBill, err := bills.GetOneBill(newTransaction.PendingBillId)
		assert.Nil(t, err)
		assert.Equal(t, newPendingBill.ParentTransactionId, newTransaction.ID)
		assert.LessOrEqual(t, newPendingBill.Amount, newTransaction.AmountWithFee)
		assert.Equal(t, newPendingBill.Amount, newTransaction.Amount) // fee 0
		assert.Equal(t, newPendingBill.Date, newTransaction.Date)
		assert.Equal(t, newPendingBill.Description, newTransaction.Description)

		// delete
		req, err := http.NewRequest(http.MethodDelete, "/transactions", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deletedTransaction := Transaction{}
		err = json.Unmarshal(w.Body.Bytes(), &deletedTransaction)
		assert.Nil(t, err)

		assert.Equal(t, newTransaction.ID, deletedTransaction.ID)
		assert.Equal(t, newTransaction.AccountId, deletedTransaction.AccountId)
		assert.Equal(t, newTransaction.Amount, deletedTransaction.Amount)
		assert.Equal(t, newTransaction.Balance, deletedTransaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt.UTC(), deletedTransaction.CreatedAt.UTC())
		assert.Equal(t, newTransaction.UpdatedAt.UTC(), deletedTransaction.UpdatedAt.UTC())
		assert.Equal(t, newTransaction.Date.UTC(), deletedTransaction.Date.UTC())
		assert.Equal(t, newTransaction.Description, deletedTransaction.Description)
		assert.Equal(t, newTransaction.PersonId, deletedTransaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, deletedTransaction.PersonName)

		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		// assert.Equal(t, deletedLastTransaction.PendingBillId, uuid.UUID{})
		assert.Equal(t, deletedTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, deletedTransaction.RevertBillId, uuid.UUID{})

		url := fmt.Sprintf("/transaction/%v", newTransaction.ID.String())
		req2, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w2.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)

		// pending bill also deleted
		_, err = bills.GetOneBill(newTransaction.PendingBillId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
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
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("Error when deleting pending bill associated with transaction", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Fee = utility.GetRandomFee()
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)
		// this deletion should be forbidden
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/pending_bills/"+newTransaction.PendingBillId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.BL003, errResponse.Error)
		assert.Equal(t, "BL003", errResponse.Code)

		// Get transaction
		sameTransaction, err := GetTransaction(newTransaction.ID)
		assert.Nil(t, err)

		assert.Equal(t, newTransaction.ID, sameTransaction.ID)
		assert.Equal(t, newTransaction.AccountId, sameTransaction.AccountId)
		assert.Equal(t, newTransaction.Amount, sameTransaction.Amount)
		assert.Equal(t, newTransaction.Balance, sameTransaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt, sameTransaction.CreatedAt)
		assert.Equal(t, newTransaction.UpdatedAt, sameTransaction.UpdatedAt)
		assert.Equal(t, newTransaction.Date, sameTransaction.Date)
		assert.Equal(t, newTransaction.Description, sameTransaction.Description)
		assert.Equal(t, newTransaction.PersonId, sameTransaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, sameTransaction.PersonName)
		assert.Equal(t, newTransaction.Currency, account.Currency)
		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		// assert.Equal(t, sameTransaction.PendingBillId, uuid.UUID{})
		assert.Equal(t, sameTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, sameTransaction.RevertBillId, uuid.UUID{})

		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Equal(t, config.Offset, transactions.Offset)
		assert.Equal(t, config.Limit, transactions.Limit)
		assert.Equal(t, 1, transactions.Count)
		assert.Len(t, transactions.Transactions, 1)
	})

	// at the end of all transactions services tests
	money_accounts.DeleteAllMoneyAccounts()
	persons.DeleteAllPersons()
}
