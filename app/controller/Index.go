package controller

import (
	Test0 "controller/Test"
	"core/library"
)

// CtlS 由于go没有类似于php的类继承，就只能通过伪继承来面向对象了
type CtlS struct {
}

// GetCtlS 获取控制器结构体映射
func (r *CtlS) GetCtlS(SS library.ServerS) map[string]interface{} {
	RegisterMessage := make(map[string]interface{})

	//----------需要在这里注册控制器包-----------------//
	RegisterMessage["CtlTest"] = &Test0.CtlTest{SS: SS}

	return RegisterMessage
}
