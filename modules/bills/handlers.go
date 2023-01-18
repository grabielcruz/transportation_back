package bills

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/julienschmidt/httprouter"
)

func GetPendingBillsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	person_id, err := uuid.Parse(ps.ByName("person_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	query := r.URL.Query()
	to_pay, err := strconv.ParseBool(query.Get("to_pay"))
	if err != nil {
		common.SendInvalidQueryStringError(w, err.Error())
		return
	}

	to_charge, err := strconv.ParseBool(query.Get("to_charge"))
	if err != nil {
		common.SendInvalidQueryStringError(w, err.Error())
		return
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		common.SendInvalidQueryStringError(w, err.Error())
		return
	}
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		common.SendInvalidQueryStringError(w, err.Error())
		return
	}
	billResponse, err := GetPendingBills(person_id, to_pay, to_charge, limit, offset)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, billResponse)
}

func CreatePendingBillHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	billFields := BillFields{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.SendReadError(w)
		return
	}
	err = json.Unmarshal(body, &billFields)
	if err != nil {
		common.SendUnmarshalError(w)
		return
	}
	// check bill fields
	err = checkBillFields(billFields)
	if err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	newBill, err := CreatePendingBill(billFields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusCreated, newBill)
}

func GetOneBillHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bill_id, err := uuid.Parse(ps.ByName("bill_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	bill, err := GetOneBill(bill_id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, bill)
}

func UpdatePendingBillHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bill_id, err := uuid.Parse(ps.ByName("bill_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	billFields := BillFields{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.SendReadError(w)
		return
	}
	err = json.Unmarshal(body, &billFields)
	if err != nil {
		common.SendUnmarshalError(w)
		return
	}
	err = checkBillFields(billFields)
	if err != nil {
		common.SendValidationError(w, err.Error())
		return
	}
	updatedBill, err := UpdatePendingBill(bill_id, billFields)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, updatedBill)
}

func DeleteBillHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bill_id, err := uuid.Parse(ps.ByName("bill_id"))
	if err != nil {
		common.SendInvalidUUIDError(w, err.Error())
		return
	}
	deletedId, err := DeleteBill(bill_id)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, deletedId)
}
