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

func GetTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	transaction_id, err := uuid.Parse(ps.ByName("transaction_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	transaction, err := GetTransaction(transaction_id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, transaction)
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	transaction := Transaction{}
	person_id, err := uuid.Parse(ps.ByName("person_id"))
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
	transaction, err = CreateTransaction(fields, person_id, true)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusCreated, transaction)
}

func DeleteLastTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	trashedTransaction, err := DeleteLastTransaction()
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, trashedTransaction)
}
