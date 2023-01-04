package persons

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
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestPersonsHandlers(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	router := httprouter.New()
	Routes(router)

	// zero person should be counted
	t.Run("Get empty slice of persons initially", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/persons", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		persons := []Person{}
		err = json.Unmarshal(w.Body.Bytes(), &persons)
		assert.Nil(t, err)
		assert.Len(t, persons, 1)
	})

	t.Run("Create one person", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GeneratePersonFields()
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/persons", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newPerson := Person{}
		err = json.Unmarshal(w.Body.Bytes(), &newPerson)
		assert.Nil(t, err)
		assert.Equal(t, fields.Name, newPerson.Name)
		assert.Equal(t, fields.Document, newPerson.Document)
	})

	DeleteAllPersons()

	t.Run("Create three persons and get an slice of three persons", func(t *testing.T) {
		CreatePerson(GeneratePersonFields())
		CreatePerson(GeneratePersonFields())
		CreatePerson(GeneratePersonFields())
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/persons", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		persons := []Person{}
		err = json.Unmarshal(w.Body.Bytes(), &persons)
		assert.Nil(t, err)
		assert.Len(t, persons, 3)
	})

	DeleteAllPersons()

	t.Run("Error when sending invalid json when creating person", func(t *testing.T) {
		buf := bytes.Buffer{}
		badFields := generateBadPersonFields()
		err := json.NewEncoder(&buf).Encode(badFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/persons", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, errors_handler.UM001, errResponse.Error)
		assert.Equal(t, "UM001", errResponse.Code)
	})

	t.Run("Error when sending bad fields on creating a person", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := PersonFields{}
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/persons", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &errResponse)
		assert.Nil(t, err)
		assert.NotNil(t, errResponse.Error)
		assert.Equal(t, "Name is required", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("Create one person and get it", func(t *testing.T) {
		fields := GeneratePersonFields()
		newPerson, err := CreatePerson(fields)
		assert.Nil(t, err)
		wantedId := newPerson.ID

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/persons/"+wantedId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		person := Person{}
		err = json.Unmarshal(w.Body.Bytes(), &person)
		assert.Nil(t, err)
		assert.Equal(t, fields.Name, person.Name)
		assert.Equal(t, fields.Document, person.Document)
	})

	t.Run("Get error when sending bad id", func(t *testing.T) {
		badId := utility.GetRandomString(10)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/persons/"+badId, nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 10", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Get error when sending uregistered id", func(t *testing.T) {
		wantedId := uuid.UUID{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/persons/"+wantedId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("It should create and update one person", func(t *testing.T) {
		createFields := GeneratePersonFields()
		newPerson, err := CreatePerson(createFields)
		assert.Nil(t, err)
		wantedId := newPerson.ID
		buf := bytes.Buffer{}
		updateFields := GeneratePersonFields()
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/persons/"+wantedId.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		updatedPerson := Person{}
		err = json.Unmarshal(w.Body.Bytes(), &updatedPerson)
		assert.Nil(t, err)
		assert.Equal(t, updateFields.Name, updatedPerson.Name)
		assert.Equal(t, updateFields.Document, updatedPerson.Document)
		assert.Greater(t, updatedPerson.UpdatedAt, updatedPerson.CreatedAt)
	})

	t.Run("Error when sending bad id", func(t *testing.T) {
		badId := utility.GetRandomString(10)
		w := httptest.NewRecorder()
		buf := bytes.Buffer{}
		updateFields := GeneratePersonFields()
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		req, err := http.NewRequest(http.MethodPatch, "/persons/"+badId, &buf)
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
		wantedId := uuid.UUID{}
		buf := bytes.Buffer{}
		updateFields := GeneratePersonFields()
		err := json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/persons/"+wantedId.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("Error when sending bad json on updating person", func(t *testing.T) {
		wantedId := uuid.UUID{}
		buf := bytes.Buffer{}
		badFields := generateBadPersonFields()
		err := json.NewEncoder(&buf).Encode(badFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/persons/"+wantedId.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.UM001, errResponse.Error)
		assert.Equal(t, "UM001", errResponse.Code)
	})

	t.Run("Error when sending bad fields on updating person", func(t *testing.T) {
		wantedId := uuid.UUID{}
		buf := bytes.Buffer{}
		badFields := PersonFields{}
		err := json.NewEncoder(&buf).Encode(badFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/persons/"+wantedId.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Name is required", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("It should create a person and delete it", func(t *testing.T) {
		fields := GeneratePersonFields()
		newPerson, err := CreatePerson(fields)
		assert.Nil(t, err)
		newId := newPerson.ID

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/persons/"+newId.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deletedId := common.ID{}
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &deletedId)
		assert.Nil(t, err)
		assert.Equal(t, newId, deletedId.ID)

		deletedAccount, err := GetOnePerson(newId)
		assert.Equal(t, deletedAccount.ID, uuid.UUID{})
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	DeleteAllPersons()

	t.Run("It should send error when sending bad id", func(t *testing.T) {
		newId := utility.GetRandomString(10)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/persons/"+newId, nil)
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

	t.Run("It should send error when trying to delete unexisting person", func(t *testing.T) {
		newId := uuid.UUID{}

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/persons/"+newId.String(), nil)
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

	t.Run("Error when creating two person with the same document", func(t *testing.T) {
		buf := bytes.Buffer{}
		fields := GeneratePersonFields()
		err := json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/persons", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		fields2 := GeneratePersonFields()
		fields2.Document = fields.Document
		err = json.NewEncoder(&buf).Encode(fields)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		req2, err := http.NewRequest(http.MethodPost, "/persons", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w2.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.PE001, errResponse.Error)
		assert.Equal(t, "PE001", errResponse.Code)
	})
}
