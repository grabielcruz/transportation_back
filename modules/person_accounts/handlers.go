package person_accounts

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/julienschmidt/httprouter"
)

func GetPersonAccountsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	person_id, err := uuid.Parse(ps.ByName("person_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}

	personAccounts, err := GetPersonAccounts(person_id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, personAccounts)
}

func CreatePersonAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fields := PersonAccountFields{}
	person_id, err := uuid.Parse(ps.ByName("person_id"))
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
	if err := checkPersonAccountFields(fields); err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	personAccount, err := CreatePersonAccount(person_id, fields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusCreated, personAccount)
}

func GetOnePersonAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account_id, err := uuid.Parse(ps.ByName("account_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	personAccount, err := GetOnePersonAccount(account_id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, personAccount)
}

func UpdatePersonAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fields := UpdatePersonAccountFields{}
	account_id, err := uuid.Parse(ps.ByName("account_id"))
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
	if err := checkUpdatePersonAccountFields(fields); err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	personAccount, err := UpdatePersonAccount(account_id, fields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, personAccount)
}

func DeletePersonAccountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account_id, err := uuid.Parse(ps.ByName("account_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	deletedId, err := DeletePersonAccount(account_id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, deletedId)
}
