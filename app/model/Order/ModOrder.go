package model

import (
	"core/library"
)

type ModOrder struct {
}

func (r *ModOrder) OrderList(Data map[string]string) map[string]interface{} {

	library.SetLog(Data)

	return map[string]interface{}{"code": "1", "msg": "子mod拉起成功"}
}
