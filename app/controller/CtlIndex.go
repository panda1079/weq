package controller

import (
	Order "controller/Test"
	"core/library"
)

// CtlIndex 由于go没有类似于php的类继承，就只能通过伪继承来面向对象了
type CtlIndex struct {
}

func (r *CtlIndex) Init(SS library.ServerS) map[string]interface{} {
	var RegisterMessage = make(map[string]interface{})

	//----------需要在这里注册控制器包-----------------//
	RegisterMessage["CtlTest"] = &Order.CtlTest{SS: SS}
	//RegisterMessage["CtlTest"] = &Order.CtlTest{SS: SS}
	//RegisterMessage["CtlTest"] = &Order.CtlTest{SS: SS}
	//RegisterMessage["CtlTest"] = &Order.CtlTest{SS: SS}

	return RegisterMessage
}
