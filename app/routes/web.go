package routes

import (
	"go/public"
)

func Web() map[string]map[string]string {
	var AddRe = public.Request{} //core.Request.Post()

	AddRe.Get("/order/order_list", "order_list", "Order/CtlOrder")

	//AddRe.Get("/order/order_list", "order_list", "controller/Order/CtlOrder")
	//AddRe.Get("/order/order_list_api", "order_list_api", "controller/Order/CtlOrder")
	//AddRe.Post("/order/replenish_order", "replenish_order", "controller/Order/CtlOrder")
	//AddRe.Post("/order/del_order", "del_order", "controller/Order/CtlOrder")

	return AddRe.RequestList
}
