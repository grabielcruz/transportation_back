package transactions

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Routes(router *httprouter.Router) {
	router.GET("/transactions/:account_id", GetTransactionsHandler)
	router.GET("/transaction/:transaction_id", GetTransactionHandler)

	// always should have a person id none zero uuid, otherwise it will throw an error
	router.POST("/transaction_to_pending_bill/:person_id", CreateTransactionHandler)

	router.POST("/close_pending_bill/:bill_id/:completed", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {})
	router.POST("/revert_closed_bill/:bill_id", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {})
	// router.POST("/transactions/:person_id", CreateTransactionHandler)
	// router.POST("/revert_pending_bill/:bill_id", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {})

	// router.PATCH("/transactions/:transaction_id", UpdateLastTransactionHandler)
	router.DELETE("/transactions", DeleteLastTransactionHandler)

}
