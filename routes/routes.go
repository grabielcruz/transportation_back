package routes

import (
	"fmt"
	"net/http"

	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/julienschmidt/httprouter"
)

func InitialHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}

func SetupAndGetRoutes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/ping", InitialHandler)

	money_accounts.Routes(router)

	return router
}
