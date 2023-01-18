package persons

import "github.com/julienschmidt/httprouter"

func Routes(router *httprouter.Router) {
	router.GET("/persons", GetPersonsHandler)
	router.POST("/persons", CreatePersonHandler)
	router.GET("/persons/:id", GetOnePersonHandler)
	router.PATCH("/persons/:id", UpdatePersonHandler)
	router.DELETE("/persons/:id", DeleteOnePersonHandler)
}
