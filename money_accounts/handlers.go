package money_accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMoneyAccountsHandler(c *gin.Context) {
	accounts := GetMoneyAccounts()
	c.JSON(200, accounts)
}

func CreateMoneyAccountHandler(c *gin.Context) {
	account := MoneyAccount{}
	fields := MoneyAccountFields{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := checkAccountFields(fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account = CreateMoneyAccount(fields)
	c.JSON(http.StatusOK, account)
}
