package person_accounts

import "github.com/julienschmidt/httprouter"

func Routes(router *httprouter.Router) {
	router.GET("/person_accounts/:person_id", GetPersonAccountsHandler)
	router.POST("/person_accounts/:person_id", CreatePersonAccountHandler)
	router.GET("/one_person_account/:account_id", GetOnePersonAccountHandler)
	router.PATCH("/one_person_account/:account_id", UpdatePersonAccountHandler)
	router.DELETE("/one_person_account/:account_id", DeletePersonAccountHandler)
}
