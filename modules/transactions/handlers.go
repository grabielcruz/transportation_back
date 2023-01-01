package transactions

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/julienschmidt/httprouter"
)

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	transactionResponse := TransationResponse{}
	account_id, err := uuid.Parse(ps.ByName("account_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	values := r.URL.Query()
	offset, err := strconv.Atoi(values.Get("offset"))
	if err != nil {
		common.SendInvalidQueryStringError(w, err.Error())
		return
	}
	limit, err := strconv.Atoi(values.Get("limit"))
	if err != nil {
		common.SendInvalidQueryStringError(w, err.Error())
		return
	}
	transactionResponse, err = GetTransactions(account_id, limit, offset)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, transactionResponse)
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	transaction := Transaction{}
	fields := TransactionFields{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.SendReadError(w)
		return
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		common.SendUnmarshalError(w)
		return
	}
	if err := checkTransactionFields(fields); err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	transaction, err = CreateTransaction(fields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusCreated, transaction)
}

func UpdateLastTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	updatedTransaction := Transaction{}
	transaction_id, err := uuid.Parse(ps.ByName("transaction_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	fields := TransactionFields{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.SendReadError(w)
		return
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		common.SendUnmarshalError(w)
		return
	}
	if err := checkTransactionFields(fields); err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	updatedTransaction, err = UpdateLastTransaction(transaction_id, fields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, updatedTransaction)
}

func DeleteLastTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	trashedTransaction, err := DeleteLastTransaction()
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, trashedTransaction)
}

func GetTrashedTransactionsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	transactions, err := GetTrashedTransactions()
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, transactions)
}

func RestoreTrashedTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	transaction_id, err := uuid.Parse(ps.ByName("transaction_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	restoredTransaction, err := RestoreTrashedTransaction(transaction_id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, restoredTransaction)
}

func DeleteTrashedTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	transaction_id, err := uuid.Parse(ps.ByName("transaction_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	deletedTransaction, err := DeleteTrashedTransaction(transaction_id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, deletedTransaction)
}
