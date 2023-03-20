package model

import (
	"core/library"
	Order "model/Test"
)

type ModIndex struct {
}

// Init 如果需要抽象化mod，就用这个内鬼吧，啊哈哈哈哈哈
func (r *ModIndex) Init(SS library.ServerS) map[string]interface{} {
	var RegisterMessage = make(map[string]interface{})

	//----------需要在这里注册控制器包-----------------//
	RegisterMessage["ModTest"] = &Order.ModTest{SS: SS}
	//RegisterMessage["ModTest"] = &Order.ModTest{}
	//RegisterMessage["ModTest"] = &Order.ModTest{}
	//RegisterMessage["ModTest"] = &Order.ModTest{}

	return RegisterMessage
}
