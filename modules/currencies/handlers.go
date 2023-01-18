package currencies

import (
	"net/http"

	"github.com/grabielcruz/transportation_back/common"
	"github.com/julienschmidt/httprouter"
)

func GetCurrenciesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currencies := GetCurrencies()
	common.SendJson(w, http.StatusOK, currencies)
}

func CreateCurrencyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	currency := ps.ByName("currency")
	createdCurrency, err := CreateCurrency(currency)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusCreated, createdCurrency)
}

func DeleteCurrencyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	currency := ps.ByName("currency")
	deletedCurrency, err := DeleteCurrency(currency)
	if err != nil {
		common.SendServiceError(w, err.Error())
		return
	}
	common.SendJson(w, http.StatusOK, deletedCurrency)
}
