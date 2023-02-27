package controller

import (
	controller "controller/Order"
	"fmt"
)

// CtlIndex 由于go没有类似于php的类继承，就只能通过伪继承来面向对象了
type CtlIndex struct {
}

func (r *CtlIndex) Init() map[string]interface{} {
	var RegisterMessage = make(map[string]interface{})

	RegisterMessage["CtlOrder"] = &controller.CtlOrder{}

	return RegisterMessage
}

func (r *CtlIndex) Test() {
	fmt.Print(1213212123)
}
