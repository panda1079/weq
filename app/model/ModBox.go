package model

import (
	"core/library"
	Test0 "model/Test"
)

// TestS 由于go没有类似于php的类继承，就只能通过伪继承来面向对象了
type TestS struct {
}

// GetTestS 获取测试结构体映射
func (r *TestS) GetTestS(SS library.ServerS) map[string]interface{} {
	var RegisterMessage = make(map[string]interface{})

	//----------需要在这里模组控制器包-----------------//

	// 测试注册
	RegisterMessage["TestModel"] = &Test0.TestModel{SS: SS}

	return RegisterMessage
}
