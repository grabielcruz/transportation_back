package common

import (
	"encoding/json"
	"log"
	"net/http"

	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func SendJson(w http.ResponseWriter, httpCode int, data any) {
	w.WriteHeader(httpCode)
	json_data, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json_data)
	return
}

func sendJsonError(w http.ResponseWriter, httpCode int, errorCode string, msg string) {
	w.WriteHeader(httpCode)
	errorResponse := errors_handler.ErrorResponse{Code: errorCode, Error: msg}
	json_data, err := json.Marshal(errorResponse)
	if err != nil {
		// should never happend
		log.Fatal(err)
	}
	w.Write(json_data)
	return
}
func SendReadError(w http.ResponseWriter) {
	sendJsonError(w, http.StatusBadRequest, "RE001", errors_handler.RE001)
}

func SendUnmarshalError(w http.ResponseWriter) {
	sendJsonError(w, http.StatusBadRequest, "UM001", errors_handler.UM001)
}

func SendValidationError(w http.ResponseWriter, msg string) {
	sendJsonError(w, http.StatusBadRequest, "VA001", msg)
}

func SendInvalidUUIDError(w http.ResponseWriter, msg string) {
	sendJsonError(w, http.StatusBadRequest, "UI001", msg)
}

func SendServiceError(w http.ResponseWriter, msg string) {
	// errors here may vary depending on the service
	errorCode := errors_handler.MapServiceError(msg)
	sendJsonError(w, http.StatusBadRequest, errorCode, msg)
}

func SendInvalidQueryStringError(w http.ResponseWriter, msg string) {
	sendJsonError(w, http.StatusBadRequest, "QS001", msg)
}
