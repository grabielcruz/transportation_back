package money_accounts

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/julienschmidt/httprouter"
)

func GetMoneyAccountsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accounts := GetMoneyAccounts()
	utility.SendJson(w, http.StatusOK, accounts)
}

func CreateMoneyAccountHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	account := MoneyAccount{}
	fields := MoneyAccountFields{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utility.SendReadError(w)
		return
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		utility.SendUnmarshalError(w)
		return
	}
	if err := checkAccountFields(fields); err != nil {
		utility.SendValidationError(w, err.Error())
		return
	}
	account = CreateMoneyAccount(fields)
	utility.SendJson(w, http.StatusOK, account)
}

func GetOneMoneyAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		utility.SendInvalidUUIDError(w, err.Error())
		return
	}
	account, err := GetOneMoneyAccount(id)
	if err != nil {
		utility.SendServiceError(w, err.Error())
		return
	}
	utility.SendJson(w, http.StatusOK, account)
}

func UpdateMoneyAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := MoneyAccount{}
	fields := MoneyAccountFields{}
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		utility.SendInvalidUUIDError(w, err.Error())
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utility.SendReadError(w)
		return
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		utility.SendUnmarshalError(w)
		return
	}
	if err := checkAccountFields(fields); err != nil {
		utility.SendValidationError(w, err.Error())
		return
	}
	account, err = UpdateMoneyAccount(id, fields)
	if err != nil {
		utility.SendServiceError(w, err.Error())
		return
	}
	utility.SendJson(w, http.StatusOK, account)
}
