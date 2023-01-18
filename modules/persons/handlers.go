package persons

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/julienschmidt/httprouter"
)

func GetPersonsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	persons := GetPersons()
	common.SendJson(w, http.StatusOK, persons)
}

func CreatePersonHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	person := Person{}
	fields := PersonFields{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.SendReadError(w)
		return
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		common.SendUnmarshalError(w)
		return
	}
	if err := checkPersonFields(fields); err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	person, err = CreatePerson(fields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusCreated, person)
}

func GetOnePersonHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	person, err := GetOnePerson(id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusCreated, person)
}

func UpdatePersonHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fields := PersonFields{}
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
	if err := checkPersonFields(fields); err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	person, err := UpdatePerson(id, fields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, person)
}

func DeleteOnePersonHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	deletedId, err := DeleteOnePerson(id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, deletedId)
}
