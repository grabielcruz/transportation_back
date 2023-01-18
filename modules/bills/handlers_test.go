package bills

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
	"github.com/grabielcruz/transportation_back/modules/config"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestBillsHandlers(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	router := httprouter.New()
	Routes(router)
	person1, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)
	person2, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)

	getBillsUrl := "/pending_bills/%v?to_pay=%v&to_charge=%v&limit=%d&offset=%d"
	t.Run("Get all pending bills response with zero bills", func(t *testing.T) {
		url := fmt.Sprintf(getBillsUrl, uuid.UUID{}, "true", "true", config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		billResponse := BillResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 0)
		assert.Equal(t, billResponse.Count, 0)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
	})

	t.Run("Error when sending bad person_id", func(t *testing.T) {
		url := fmt.Sprintf(getBillsUrl, "000", "true", "true", config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 3", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Error when sending bad to_pay field", func(t *testing.T) {
		url := fmt.Sprintf(getBillsUrl, uuid.UUID{}, "asalto", "true", config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "strconv.ParseBool: parsing \"asalto\": invalid syntax", errResponse.Error)
		assert.Equal(t, "QS001", errResponse.Code)
	})

	t.Run("Error when sending bad to_charge field", func(t *testing.T) {
		url := fmt.Sprintf(getBillsUrl, uuid.UUID{}, "true", "asalto", config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "strconv.ParseBool: parsing \"asalto\": invalid syntax", errResponse.Error)
		assert.Equal(t, "QS001", errResponse.Code)
	})

	t.Run("Error when sending bad limit field", func(t *testing.T) {
		url := fmt.Sprintf(getBillsUrl, uuid.UUID{}, "true", "true", "A", config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "strconv.Atoi: parsing \"\": invalid syntax", errResponse.Error)
		assert.Equal(t, "QS001", errResponse.Code)
	})

	t.Run("Error when sending bad offset field", func(t *testing.T) {
		url := fmt.Sprintf(getBillsUrl, uuid.UUID{}, "true", "true", config.Limit, "A")
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "strconv.Atoi: parsing \"\": invalid syntax", errResponse.Error)
		assert.Equal(t, "QS001", errResponse.Code)
	})

	t.Run("Create one pending bill", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/pending_bills", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newBill := Bill{}
		err = json.Unmarshal(w.Body.Bytes(), &newBill)
		assert.Nil(t, err)
		assert.Equal(t, billFields.PersonId, newBill.PersonId)
		assert.Equal(t, person1.Name, newBill.PersonName)
		assert.Equal(t, billFields.Currency, newBill.Currency)
		assert.Equal(t, billFields.Date.Format("2006-01-02"), newBill.Date.Format("2006-01-02"))
		assert.Equal(t, billFields.Description, newBill.Description)
		assert.Equal(t, billFields.Amount, newBill.Amount)
	})

	EmptyBills()

	t.Run("Error when creating bill with zero person id", func(t *testing.T) {
		billFields := GenerateBillFields(uuid.UUID{})
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/pending_bills", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Person id should be not zero uuid", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("Error when creating bill with empty description", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		billFields.Description = ""
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/pending_bills", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Description is required", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("Error when creating bill with amount zero", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		billFields.Amount = 0
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/pending_bills", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Amount should be greater than zero", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("Error when creating bill with bad currency", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		billFields.Currency = ""
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/pending_bills", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Currency code should be 3 upper case letters", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("Create 4 bills, 2 for person1, 2 for person2, negative and positive balance and get them filtered", func(t *testing.T) {
		// person1
		billFields := GenerateBillFields(person1.ID)
		billFields.Amount = 55.55
		_, err := CreatePendingBill(billFields)
		assert.Nil(t, err)

		billFields = GenerateBillFields(person1.ID)
		billFields.Amount = -55.55
		_, err = CreatePendingBill(billFields)
		assert.Nil(t, err)

		// person2
		billFields = GenerateBillFields(person2.ID)
		billFields.Amount = 77.77
		_, err = CreatePendingBill(billFields)
		assert.Nil(t, err)

		billFields = GenerateBillFields(person2.ID)
		billFields.Amount = -77.77
		_, err = CreatePendingBill(billFields)
		assert.Nil(t, err)

		// all of them
		// billResponse, err := GetPendingBills(uuid.UUID{}, true, true, config.Limit, config.Offset)
		url := fmt.Sprintf(getBillsUrl, uuid.UUID{}, "true", "true", config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		billResponse := BillResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 4)
		assert.Equal(t, billResponse.Count, 4)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person2.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[1].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[2].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[2].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[3].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[3].Amount)

		// person1
		// billResponse, err = GetPendingBills(person1.ID, true, true, config.Limit, config.Offset)
		url = fmt.Sprintf(getBillsUrl, person1.ID, "true", "true", config.Limit, config.Offset)
		req2, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		err = json.Unmarshal(w2.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 2)
		assert.Equal(t, billResponse.Count, 2)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[1].Amount)

		// person2
		// billResponse, err = GetPendingBills(person2.ID, true, true, config.Limit, config.Offset)
		url = fmt.Sprintf(getBillsUrl, person2.ID, "true", "true", config.Limit, config.Offset)
		req3, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w2.Code)

		err = json.Unmarshal(w3.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 2)
		assert.Equal(t, billResponse.Count, 2)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person2.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[1].Amount)

		// to_charge only
		// billResponse, err = GetPendingBills(uuid.UUID{}, false, true, config.Limit, config.Offset)
		url = fmt.Sprintf(getBillsUrl, uuid.UUID{}, "false", "true", config.Limit, config.Offset)
		req4, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w4 := httptest.NewRecorder()
		router.ServeHTTP(w4, req4)
		assert.Equal(t, http.StatusOK, w4.Code)

		err = json.Unmarshal(w4.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 2)
		assert.Equal(t, billResponse.Count, 2)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, billResponse.Bills[0].PersonId, person2.ID)
		assert.Equal(t, float64(77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[1].Amount)

		// to_pay only
		// billResponse, err = GetPendingBills(uuid.UUID{}, true, false, config.Limit, config.Offset)
		url = fmt.Sprintf(getBillsUrl, uuid.UUID{}, "true", "false", config.Limit, config.Offset)
		req5, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w5 := httptest.NewRecorder()
		router.ServeHTTP(w5, req5)
		assert.Equal(t, http.StatusOK, w5.Code)

		err = json.Unmarshal(w5.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 2)
		assert.Equal(t, billResponse.Count, 2)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[1].Amount)

		// person1 to_charge
		// billResponse, err = GetPendingBills(person1.ID, false, true, config.Limit, config.Offset)
		url = fmt.Sprintf(getBillsUrl, person1.ID, "false", "true", config.Limit, config.Offset)
		req6, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w6 := httptest.NewRecorder()
		router.ServeHTTP(w6, req6)
		assert.Equal(t, http.StatusOK, w6.Code)

		err = json.Unmarshal(w6.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 1)
		assert.Equal(t, billResponse.Count, 1)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[0].Amount)

		// person1 to_pay
		// billResponse, err = GetPendingBills(person1.ID, true, false, config.Limit, config.Offset)
		url = fmt.Sprintf(getBillsUrl, person1.ID, "true", "false", config.Limit, config.Offset)
		req7, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w7 := httptest.NewRecorder()
		router.ServeHTTP(w7, req7)
		assert.Equal(t, http.StatusOK, w7.Code)

		err = json.Unmarshal(w7.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 1)
		assert.Equal(t, billResponse.Count, 1)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[0].Amount)

		// person2 to_charge
		// billResponse, err = GetPendingBills(person2.ID, false, true, config.Limit, config.Offset)
		url = fmt.Sprintf(getBillsUrl, person2.ID, "false", "true", config.Limit, config.Offset)
		req8, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w8 := httptest.NewRecorder()
		router.ServeHTTP(w8, req8)
		assert.Equal(t, http.StatusOK, w8.Code)

		err = json.Unmarshal(w8.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 1)
		assert.Equal(t, billResponse.Count, 1)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[0].Amount)

		// person2 to_pay
		// billResponse, err = GetPendingBills(person2.ID, true, false, config.Limit, config.Offset)
		url = fmt.Sprintf(getBillsUrl, person2.ID, "true", "false", config.Limit, config.Offset)
		req9, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w9 := httptest.NewRecorder()
		router.ServeHTTP(w9, req9)
		assert.Equal(t, http.StatusOK, w9.Code)

		err = json.Unmarshal(w9.Body.Bytes(), &billResponse)
		assert.Nil(t, err)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 1)
		assert.Equal(t, billResponse.Count, 1)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
	})

	EmptyBills()

	t.Run("Error when requesting not to pay and not to charge", func(t *testing.T) {
		url := fmt.Sprintf(getBillsUrl, uuid.UUID{}, "false", "false", config.Limit, config.Offset)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.BL001, errResponse.Error)
		assert.Equal(t, "BL001", errResponse.Code)
	})

	t.Run("Error when creating a bill with balance less than zero", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		billFields.Amount = -55
		err := json.NewEncoder(&buf).Encode(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/pending_bills", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Amount should be greater than zero", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	t.Run("Create one pending bill and get it with single response", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/pending_bills", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newBill := Bill{}
		err = json.Unmarshal(w.Body.Bytes(), &newBill)
		assert.Nil(t, err)

		w2 := httptest.NewRecorder()
		req2, err := http.NewRequest(http.MethodGet, "/bills/"+newBill.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		gotBill := Bill{}
		err = json.Unmarshal(w2.Body.Bytes(), &gotBill)
		assert.Nil(t, err)
		assert.Equal(t, newBill.ID, gotBill.ID)
		assert.Equal(t, newBill.PersonId, gotBill.PersonId)
		assert.Equal(t, newBill.PersonName, gotBill.PersonName)
		assert.Equal(t, newBill.Date.Format("2006-01-02"), gotBill.Date.Format("2006-01-02"))
		assert.Equal(t, newBill.Description, gotBill.Description)
		assert.Equal(t, newBill.Currency, gotBill.Currency)
		assert.Equal(t, newBill.Amount, gotBill.Amount)
		assert.Equal(t, newBill.CreatedAt, gotBill.CreatedAt)
		assert.Equal(t, newBill.UpdatedAt, gotBill.UpdatedAt)
	})

	EmptyBills()

	t.Run("Create closed bill artifitially and get it with single response", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		newBill, err := createClosedBill(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/bills/"+newBill.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		gotBill := Bill{}
		err = json.Unmarshal(w.Body.Bytes(), &gotBill)
		assert.Nil(t, err)
		assert.Equal(t, newBill.ID, gotBill.ID)
		assert.Equal(t, newBill.PersonId, gotBill.PersonId)
		assert.Equal(t, newBill.PersonName, gotBill.PersonName)
		assert.Equal(t, newBill.Date.Format("2006-01-02"), gotBill.Date.Format("2006-01-02"))
		assert.Equal(t, newBill.Description, gotBill.Description)
		assert.Equal(t, newBill.Currency, gotBill.Currency)
		assert.Equal(t, newBill.Amount, gotBill.Amount)
		assert.Equal(t, newBill.CreatedAt.UTC(), gotBill.CreatedAt.UTC())
		assert.Equal(t, newBill.UpdatedAt.UTC(), gotBill.UpdatedAt.UTC())
	})

	EmptyBills()

	t.Run("Error when requesting unexisting bill", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/bills/"+(uuid.UUID{}).String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("Error when requesting bad bill id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/bills/"+"1234", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 4", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Create one bill and update it", func(t *testing.T) {
		// create bill
		billFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(billFields)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/pending_bills", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		newBill := Bill{}
		err = json.Unmarshal(w.Body.Bytes(), &newBill)
		assert.Nil(t, err)

		// update bill
		updateFields := GenerateBillFields(person1.ID)
		buf2 := bytes.Buffer{}
		err = json.NewEncoder(&buf2).Encode(updateFields)
		assert.Nil(t, err)
		w2 := httptest.NewRecorder()
		req2, err := http.NewRequest(http.MethodPatch, "/pending_bills/"+newBill.ID.String(), &buf2)
		assert.Nil(t, err)

		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		updatedBill := Bill{}
		err = json.Unmarshal(w2.Body.Bytes(), &updatedBill)
		assert.Nil(t, err)

		assert.Equal(t, updatedBill.ID, newBill.ID)
		assert.Equal(t, updatedBill.PersonId, updateFields.PersonId)
		assert.Equal(t, updatedBill.Date.Format("2006-01-02"), updateFields.Date.Format("2006-01-02"))
		assert.Equal(t, updatedBill.Description, updateFields.Description)
		assert.Equal(t, updatedBill.Currency, updateFields.Currency)
		assert.Equal(t, updatedBill.Amount, updateFields.Amount)

		// get updated bill
		w3 := httptest.NewRecorder()
		req3, err := http.NewRequest(http.MethodGet, "/bills/"+newBill.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w3.Code)

		gotBill := Bill{}
		err = json.Unmarshal(w3.Body.Bytes(), &gotBill)
		assert.Nil(t, err)

		assert.Equal(t, updatedBill.ID, gotBill.ID)
		assert.Equal(t, updatedBill.PersonId, gotBill.PersonId)
		assert.Equal(t, updatedBill.PersonName, gotBill.PersonName)
		assert.Equal(t, updatedBill.Date.Format("2006-01-02"), gotBill.Date.Format("2006-01-02"))
		assert.Equal(t, updatedBill.Description, gotBill.Description)
		assert.Equal(t, updatedBill.Currency, gotBill.Currency)
		assert.Equal(t, updatedBill.Amount, gotBill.Amount)
		assert.Equal(t, updatedBill.CreatedAt.UTC(), gotBill.CreatedAt.UTC())
		assert.Equal(t, updatedBill.UpdatedAt.UTC(), gotBill.UpdatedAt.UTC())
	})

	EmptyBills()

	t.Run("Error when updating unexisting bill", func(t *testing.T) {
		updateFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/pending_bills/"+uuid.UUID{}.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("Error when updating with bad bill id", func(t *testing.T) {
		updateFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/pending_bills/"+"1234", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 4", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})

	t.Run("Error when updating with zero person ID", func(t *testing.T) {
		bill, err := CreatePendingBill(GenerateBillFields(person1.ID))
		assert.Nil(t, err)

		updateFields := GenerateBillFields(uuid.UUID{})
		buf := bytes.Buffer{}
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/pending_bills/"+bill.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Person id should be not zero uuid", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	EmptyBills()

	t.Run("Error when updating with empty description", func(t *testing.T) {
		bill, err := CreatePendingBill(GenerateBillFields(person1.ID))
		assert.Nil(t, err)

		updateFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		updateFields.Description = ""
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/pending_bills/"+bill.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Description is required", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	EmptyBills()

	t.Run("Error when updating with negative amount", func(t *testing.T) {
		bill, err := CreatePendingBill(GenerateBillFields(person1.ID))
		assert.Nil(t, err)

		updateFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		updateFields.Amount = 0
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/pending_bills/"+bill.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Amount should be greater than zero", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	EmptyBills()

	t.Run("Error when updating with invalid currency", func(t *testing.T) {
		bill, err := CreatePendingBill(GenerateBillFields(person1.ID))
		assert.Nil(t, err)

		updateFields := GenerateBillFields(person1.ID)
		buf := bytes.Buffer{}
		updateFields.Currency = "LuisCruz"
		err = json.NewEncoder(&buf).Encode(updateFields)
		assert.Nil(t, err)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/pending_bills/"+bill.ID.String(), &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "Currency code should be 3 upper case letters", errResponse.Error)
		assert.Equal(t, "VA001", errResponse.Code)
	})

	EmptyBills()

	t.Run("Create and delete one bill", func(t *testing.T) {
		bill, err := CreatePendingBill(GenerateBillFields(person1.ID))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/pending_bills/"+bill.ID.String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		deleted_id := common.ID{}
		err = json.Unmarshal(w.Body.Bytes(), &deleted_id)
		assert.Nil(t, err)
		assert.Equal(t, deleted_id.ID, bill.ID)

		_, err = GetOneBill(bill.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Error when deleting unexisting id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/pending_bills/"+(uuid.UUID{}).String(), nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.DB001, errResponse.Error)
		assert.Equal(t, "DB001", errResponse.Code)
	})

	t.Run("Error when deleting with bad id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/pending_bills/"+"12345", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Nil(t, err)
		assert.Equal(t, "invalid UUID length: 5", errResponse.Error)
		assert.Equal(t, "UI001", errResponse.Code)
	})
}
