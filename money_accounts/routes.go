package money_accounts

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	r.GET("/money_accounts", GetMoneyAccountsHandler)
	r.POST("/money_accounts", CreateMoneyAccountHandler)
}
