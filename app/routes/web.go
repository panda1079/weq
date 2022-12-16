package routes

import (
	"core"
	"fmt"
)

func Web() {
	var AddRe = core.Request{} //core.Request.Post()

	AddRe.Get("/order/order_list", "order_list", "controller/Order/CtlOrder")
	AddRe.Get("/order/order_list_api", "order_list_api", "controller/Order/CtlOrder")
	AddRe.Post("/order/replenish_order", "replenish_order", "controller/Order/CtlOrder")
	AddRe.Post("/order/del_order", "del_order", "controller/Order/CtlOrder")

	fmt.Print(AddRe.GetList)
	fmt.Print(AddRe.PostList)
}
