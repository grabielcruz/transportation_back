package transactions

import (
	"net/http"

	"github.com/grabielcruz/transportation_back/common"
	"github.com/julienschmidt/httprouter"
)

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transactionResponse := TransationResponse{}
	common.SendJson(w, http.StatusOK, transactionResponse)
}
