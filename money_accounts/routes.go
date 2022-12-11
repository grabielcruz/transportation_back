package money_accounts

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	r.GET("/money_accounts", GetMoneyAccountsHandler)
	r.GET("/money_accounts/:id", GetOneMoneyAccountHandler)
	r.POST("/money_accounts", CreateMoneyAccountHandler)
	r.PATCH("/money_accounts/:id", UpdateMoneyAccountHandler)
}
