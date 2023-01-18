package currencies

import "github.com/julienschmidt/httprouter"

func Routes(router *httprouter.Router) {
	router.GET("/currencies", GetCurrenciesHandler)
	router.POST("/currencies/:currency", CreateCurrencyHandler)
	router.DELETE("/currencies/:currency", DeleteCurrencyHandler)
}
