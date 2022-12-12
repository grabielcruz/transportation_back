package common

import (
	"encoding/json"
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

func sendJsonError(w http.ResponseWriter, httpCode int, msg string) {
	w.WriteHeader(httpCode)
	errorResponse := errors_handler.ErrorResponse{Error: msg}
	json_data, err := json.Marshal(errorResponse)
	if err != nil {
		// should never happend
		errors_handler.CheckError(err)
	}
	w.Write(json_data)
	return
}
func SendReadError(w http.ResponseWriter) {
	sendJsonError(w, http.StatusBadRequest, "Unable to read body of the request")
}

func SendUnmarshalError(w http.ResponseWriter) {
	sendJsonError(w, http.StatusBadRequest, "Invalid data type")
}

func SendValidationError(w http.ResponseWriter, msg string) {
	sendJsonError(w, http.StatusBadRequest, msg)
}

func SendInvalidUUIDError(w http.ResponseWriter, msg string) {
	sendJsonError(w, http.StatusBadRequest, msg)
}

func SendServiceError(w http.ResponseWriter, msg string) {
	sendJsonError(w, http.StatusBadRequest, msg)
}
