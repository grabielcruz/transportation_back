package money_accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetMoneyAccountsHandler(c *gin.Context) {
	accounts := GetMoneyAccounts()
	c.JSON(200, accounts)
}

func CreateMoneyAccountHandler(c *gin.Context) {
	account := MoneyAccount{}
	fields := MoneyAccountFields{}
	if err := c.BindJSON(&fields); err != nil {
		// c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// return
	}
	if err := checkAccountFields(fields); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account = CreateMoneyAccount(fields)
	c.JSON(http.StatusOK, account)
}

func GetOneMoneyAccountHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := GetOneMoneyAccount(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func UpdateMoneyAccountHandler(c *gin.Context) {
	account := MoneyAccount{}
	fields := MoneyAccountFields{}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.BindJSON(&fields); err != nil {
		// c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// return
	}
	// if err := checkAccountFields(fields); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	account, err = UpdateMoneyAccount(id, fields)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}
