package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitialHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", InitialHandler)
	return r
}
