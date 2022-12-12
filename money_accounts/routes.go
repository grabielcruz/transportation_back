package money_accounts

import (
	"github.com/julienschmidt/httprouter"
)

func Routes(router *httprouter.Router) {
	router.GET("/money_accounts", GetMoneyAccountsHandler)
	router.GET("/money_accounts/:id", GetOneMoneyAccountHandler)
	router.POST("/money_accounts", CreateMoneyAccountHandler)
	router.PATCH("/money_accounts/:id", UpdateMoneyAccountHandler)
}
