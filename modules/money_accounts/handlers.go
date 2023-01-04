package money_accounts

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/julienschmidt/httprouter"
)

func GetMoneyAccountsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accounts := GetMoneyAccounts()
	common.SendJson(w, http.StatusOK, accounts)
}

func CreateMoneyAccountHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	account := MoneyAccount{}
	fields := MoneyAccountFields{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.SendReadError(w)
		return
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		common.SendUnmarshalError(w)
		return
	}
	if err := checkAccountFields(fields); err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	account, err = CreateMoneyAccount(fields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusCreated, account)
}

func GetOneMoneyAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	account, err := GetOneMoneyAccount(id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, account)
}

func UpdateMoneyAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fields := MoneyAccountFields{}
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.SendReadError(w)
		return
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		common.SendUnmarshalError(w)
		return
	}
	if err := checkAccountFields(fields); err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	account, err := UpdateMoneyAccount(id, fields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, account)
}

func DeleteOneMoneyAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	deletedId, err := DeleteOneMoneyAccount(id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, deletedId)
}
