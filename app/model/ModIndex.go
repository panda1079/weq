package model

import (
	Order "model/Order"
)

type ModIndex struct {
}

// Init 如果需要抽象化mod，就用这个内鬼吧，啊哈哈哈哈哈
func (r *ModIndex) Init() map[string]interface{} {
	var RegisterMessage = make(map[string]interface{})

	//----------需要在这里注册控制器包-----------------//
	RegisterMessage["ModOrder"] = &Order.ModOrder{}
	//RegisterMessage["ModOrder"] = &Order.CtlOrder{}
	//RegisterMessage["ModOrder"] = &Order.CtlOrder{}
	//RegisterMessage["ModOrder"] = &Order.CtlOrder{}

	return RegisterMessage
}
