package person_accounts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestPersonAccountHandler(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	router := httprouter.New()
	Routes(router)
	person, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)

	t.Run("Get empty slice of person accounts initially", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/person_accounts/"+person.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		personAccounts := []PersonAccount{}
		err = json.Unmarshal(w.Body.Bytes(), &personAccounts)
		assert.Nil(t, err)
		assert.Len(t, personAccounts, 0)
	})

	t.Run("Create a person account", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GeneratePersonAccountFields()
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/person_accounts/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		personAccount := PersonAccount{}
		err = json.Unmarshal(w.Body.Bytes(), &personAccount)
		assert.Nil(t, err)
		assert.Equal(t, fields.Name, personAccount.Name)
		assert.Equal(t, fields.Description, personAccount.Description)
		assert.Equal(t, fields.Currency, personAccount.Currency)
	})

	DeleteAllPersonAccounts()

	t.Run("Create three person accounts and get an slice of three person accounts", func(t *testing.T) {
		CreatePersonAccount(person.ID, GeneratePersonAccountFields())
		CreatePersonAccount(person.ID, GeneratePersonAccountFields())
		CreatePersonAccount(person.ID, GeneratePersonAccountFields())
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/person_accounts/"+person.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		personAccounts := []PersonAccount{}
		err = json.Unmarshal(w.Body.Bytes(), &personAccounts)
		assert.Nil(t, err)
		assert.Len(t, personAccounts, 3)
	})

	DeleteAllPersonAccounts()

	t.Run("Error when sending invalid json when creating a person account", func(t *testing.T) {
		buf := bytes.Buffer{}
		badFields := generateBadPersonAccountFields()
		err := json.NewEncoder(&buf).Encode(badFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/person_accounts/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.UM001, errResponse.Error)
		assert.Equal(t, "UM001", errResponse.Code)
	})

	t.Run("Testing fields validator", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GeneratePersonAccountFields()
		fields.Description = ""
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/person_accounts/"+person.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Description is required", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("Create one person account and get it", func(t *testing.T) {
		fields := GeneratePersonAccountFields()
		createdAccount, err := CreatePersonAccount(person.ID, fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/one_person_account/"+createdAccount.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		obtainedAccount := PersonAccount{}
		err = json.Unmarshal(w.Body.Bytes(), &obtainedAccount)
		assert.Nil(t, err)
		assert.Equal(t, createdAccount.ID, obtainedAccount.ID)
		assert.Equal(t, createdAccount.Name, obtainedAccount.Name)
		assert.Equal(t, createdAccount.Description, obtainedAccount.Description)
		assert.Equal(t, createdAccount.CreatedAt.UTC(), obtainedAccount.CreatedAt.UTC())
		assert.Equal(t, createdAccount.UpdatedAt.UTC(), obtainedAccount.UpdatedAt.UTC())
	})

	DeleteAllPersonAccounts()

	t.Run("Error when requesting a person account with a bad id", func(t *testing.T) {
		badId := utility.GetRandomString(10)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/one_person_account/"+badId, nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 10", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Error when person account does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/one_person_account/"+(uuid.UUID{}).String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("It should create an update one person account", func(t *testing.T) {
		fields := GeneratePersonAccountFields()
		createdAccount, err := CreatePersonAccount(person.ID, fields)
		assert.Nil(t, err)

		buf := bytes.Buffer{}
		updateFields := GeneratePersonAccountFields()
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/one_person_account/"+createdAccount.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		updatedAccount := PersonAccount{}
		err = json.Unmarshal(w.Body.Bytes(), &updatedAccount)
		assert.Nil(t, err)
		assert.Equal(t, createdAccount.ID, updatedAccount.ID)
		assert.Equal(t, updatedAccount.Name, updateFields.Name)
		assert.Equal(t, updatedAccount.Description, updateFields.Description)
		assert.Equal(t, updatedAccount.Currency, createdAccount.Currency)
	})

	DeleteAllPersonAccounts()

	t.Run("Error when sending bad id on patch", func(t *testing.T) {
		badId := utility.GetRandomString(10)
		buf := bytes.Buffer{}
		updateFields := GeneratePersonAccountFields()
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/one_person_account/"+badId, &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 10", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Error when sending unregistered id when patching", func(t *testing.T) {
		buf := bytes.Buffer{}
		updateFields := GeneratePersonAccountFields()
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/one_person_account/"+(uuid.UUID{}).String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("Error when sending bad json when patching", func(t *testing.T) {
		buf := bytes.Buffer{}
		updateFields := generateBadPersonAccountFields()
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/one_person_account/"+(uuid.UUID{}).String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.UM001, errResponse.Error)
		assert.Equal(t, "UM001", errResponse.Code)
	})

	t.Run("Error when sending bad fields when patching", func(t *testing.T) {
		buf := bytes.Buffer{}
		updateFields := GeneratePersonAccountFields()
		updateFields.Name = ""
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/one_person_account/"+(uuid.UUID{}).String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Name is required", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("It should create one person account and delete it", func(t *testing.T) {
		fields := GeneratePersonAccountFields()
		newAccount, err := CreatePersonAccount(person.ID, fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/one_person_account/"+newAccount.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deletedId := common.ID{}
		err = json.Unmarshal(w.Body.Bytes(), &deletedId)
		assert.Nil(t, err)
		assert.Equal(t, newAccount.ID, deletedId.ID)

		_, err = GetOnePersonAccount(deletedId.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	DeleteAllPersonAccounts()

	t.Run("Error when sending bad id when deleting", func(t *testing.T) {
		newId := utility.GetRandomString(10)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/one_person_account/"+newId, nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, "invalid UUID length: 10", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Error when trying to delete unexisting account", func(t *testing.T) {
		newId := uuid.UUID{}.String()
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/one_person_account/"+newId, nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

}
