package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

const (
	hola = "Hola a todos"
	que  = "Que haceis"
	hace = "vos?"
)

func jsonHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dummy := utility.DummyType{
		Hola: hola,
		Que:  que,
		Hace: hace,
	}
	SendJson(w, http.StatusOK, dummy)
}

func readBodyToJasonHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var dummy utility.DummyType
	read, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendReadError(w)
		return
	}
	err = json.Unmarshal(read, &dummy)
	if err != nil {
		SendUnmarshalError(w)
		return
	}
	SendJson(w, http.StatusOK, dummy)
}

func TestSendJson(t *testing.T) {
	router := httprouter.New()

	router.GET("/json_handler", jsonHandler)
	router.POST("/body_to_json_handler", readBodyToJasonHandler)

	t.Run("It should send a json dummy data type with status ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/json_handler", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		dummy := utility.DummyType{}
		err = json.Unmarshal(w.Body.Bytes(), &dummy)
		assert.Nil(t, err)
		assert.Equal(t, dummy.Hola, hola)
		assert.Equal(t, dummy.Que, que)
		assert.Equal(t, dummy.Hace, hace)
	})

	t.Run("It should send a json dummy data a receive it in return", func(t *testing.T) {
		var buf bytes.Buffer
		dummy := utility.GenerateDummyData()
		err := json.NewEncoder(&buf).Encode(dummy)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/body_to_json_handler", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		createdDummy := utility.DummyType{}
		err = json.Unmarshal(w.Body.Bytes(), &createdDummy)
		assert.Nil(t, err)
		assert.Equal(t, dummy.Hola, createdDummy.Hola)
		assert.Equal(t, dummy.Que, createdDummy.Que)
		assert.Equal(t, dummy.Hace, createdDummy.Hace)
	})

	t.Run("It should get error when sending wrong data type", func(t *testing.T) {
		var buf bytes.Buffer
		dummy := struct {
			Hola bool
			Que  bool
			Hace bool
		}{
			Hola: false,
			Que:  false,
			Hace: false,
		}

		err := json.NewEncoder(&buf).Encode(dummy)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/body_to_json_handler", &buf)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		errorResponse := errors_handler.ErrorResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.Nil(t, err)
		assert.Equal(t, errors_handler.UM001, errorResponse.Error)

	})
}
