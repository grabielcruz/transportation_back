package transactions

import "github.com/julienschmidt/httprouter"

func Routes(router *httprouter.Router) {
	router.GET("/transactions/:account_id", GetTransactionsHandler)
	router.POST("/transactions", CreateTransactionHandler)
	router.PATCH("/transactions/:transaction_id", UpdateLastTransactionHandler)
	router.DELETE("/transactions", DeleteLastTransactionHandler)

	router.GET("/trashed_transactions", GetTrashedTransactionsHandler)
	router.POST("/trashed_transactions/:transaction_id", RestoreTrashedTransactionHandler)
	router.DELETE("/trashed_transactions/:transaction_id", DeleteTrashedTransactionHandler)
}
