package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grabielcruz/transportation_back/money_accounts"
)

func InitialHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func SetupAndGetRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", InitialHandler)

	money_accounts.Routes(r)

	return r
}
