package transactions

import "github.com/julienschmidt/httprouter"

func Routes(router *httprouter.Router) {
	router.GET("/transactions/:account_id", GetTransactionsHandler)
	router.GET("/transaction/:transaction_id", GetTransactionHandler)
	router.POST("/transactions", CreateTransactionHandler)
	router.PATCH("/transactions/:transaction_id", UpdateLastTransactionHandler)
	router.DELETE("/transactions", DeleteLastTransactionHandler)

}
