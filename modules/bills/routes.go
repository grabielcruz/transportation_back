package bills

import "github.com/julienschmidt/httprouter"

func Routes(router *httprouter.Router) {
	router.GET("/pending_bills/:person_id", GetPendingBillsHandler)
	router.POST("/pending_bills", CreatePendingBillHandler)
	router.GET("/bills/:bill_id", GetOneBillHandler)
	router.PATCH("/pending_bills/:bill_id", UpdatePendingBillHandler)
	router.DELETE("/pending_bills/:bill_id", DeleteBillHandler)
}
