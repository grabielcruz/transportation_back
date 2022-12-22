package transactions

import "github.com/julienschmidt/httprouter"

func Routes(router *httprouter.Router) {
	router.GET("/transactions/:account_id", GetTransactionsHandler)
}
