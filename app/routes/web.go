package routes

import (
	"core/library"
)

func Web() map[string]map[string]string {
	var AddRe = library.Request{} //core.Request.Post()

	AddRe.Get("/order/order_list", "OrderList", "CtlOrder")
	AddRe.Post("/test", "TestA", "CtlOrder")

	//AddRe.Get("/order/order_list", "order_list", "controller/Order/CtlOrder")
	//AddRe.Get("/order/order_list_api", "order_list_api", "controller/Order/CtlOrder")
	//AddRe.Post("/order/replenish_order", "replenish_order", "controller/Order/CtlOrder")
	//AddRe.Post("/order/del_order", "del_order", "controller/Order/CtlOrder")

	return AddRe.RequestList
}
