package routes

import (
	"fmt"
	"net/http"

	"github.com/grabielcruz/transportation_back/modules/bills"
	"github.com/grabielcruz/transportation_back/modules/currencies"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/grabielcruz/transportation_back/modules/person_accounts"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/modules/transactions"
	"github.com/julienschmidt/httprouter"
)

func InitialHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}

func SetupAndGetRoutes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/ping", InitialHandler)

	currencies.Routes(router)
	money_accounts.Routes(router)
	persons.Routes(router)
	person_accounts.Routes(router)
	bills.Routes(router)
	transactions.Routes(router)

	return router
}
